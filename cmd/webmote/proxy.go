package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/coder/websocket"
)

var id atomic.Uint64

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	pid := id.Add(1)
	slog.Info("New proxy connection", slog.Uint64("id", pid), slog.String("source", r.RemoteAddr))

	runProxy(w, r)

	slog.Info("Closing proxy connection", slog.Uint64("id", pid), slog.String("source", r.RemoteAddr))
}

func runProxy(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("Failed to accept proxy socket:", slog.String("reason", err.Error()))
		return
	}
	defer ws.Close(websocket.StatusNormalClosure, "")

	prx := &proxy{ctx: r.Context(), ws: ws}
	err = prx.connect()
	if err != nil {
		slog.Error("Failed to connect to amplifier:", slog.String("reason", err.Error()))
		return
	}
	defer prx.amp.Close()

	go prx.forwardFromAmplifier()
	prx.forwardFromClient()
}

type proxy struct {
	ctx context.Context
	ws  *websocket.Conn
	amp net.Conn
}

func (p *proxy) connect() error {
	_, host, err := p.ws.Read(p.ctx)
	if err != nil {
		return err
	}

	p.amp, err = net.Dial("tcp", string(host)+":50001")
	return err
}

func (p *proxy) forwardFromAmplifier() {
	buf := make([]byte, 32)
	for {
		n, err := p.amp.Read(buf)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				slog.Error("Error reading from amplifier", slog.String("reason", err.Error()))
			}
			return
		}

		err = p.ws.Write(p.ctx, websocket.MessageText, buf[:n])
		if err != nil {
			slog.Error("Error writing to socket", slog.String("reason", err.Error()))
			return
		}
	}
}

func (p *proxy) forwardFromClient() {
	for {
		_, data, err := p.ws.Read(p.ctx)
		if err != nil {
			if !isAcceptedError(err) {
				slog.Error("Error reading from socket", slog.String("reason", err.Error()))
			}
			return
		}

		_, err = p.amp.Write(data)
		if err != nil {
			slog.Error("Error writing to amplifier", slog.String("reason", err.Error()))
			return
		}
	}
}

func isAcceptedError(err error) bool {
	status := websocket.CloseStatus(err)
	return status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway
}
