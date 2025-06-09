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

	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("Failed to accept proxy socket:", slog.String("reason", err.Error()))
		return
	}
	defer ws.Close(websocket.StatusNormalClosure, "")

	amp, err := connect(ws)
	if err != nil {
		slog.Error("Failed to connect to amplifier:", slog.String("reason", err.Error()))
		return
	}
	defer amp.Close()

	go forwardFromAmplifier(amp, ws)
	forwardFromClient(amp, ws)

	slog.Info("Closing proxy connection", slog.Uint64("id", pid), slog.String("source", r.RemoteAddr))
}

func connect(ws *websocket.Conn) (net.Conn, error) {
	_, host, err := ws.Read(context.Background())
	if err != nil {
		return nil, err
	}

	return net.Dial("tcp", string(host)+":50001")
}

func forwardFromAmplifier(amp net.Conn, ws *websocket.Conn) {
	buf := make([]byte, 32)
	for {
		n, err := amp.Read(buf)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				slog.Error("Error reading from amplifier", slog.String("reason", err.Error()))
			}
			return
		}

		err = ws.Write(context.Background(), websocket.MessageText, buf[:n])
		if err != nil {
			slog.Error("Error writing to socket", slog.String("reason", err.Error()))
			return
		}
	}
}

func forwardFromClient(amp net.Conn, ws *websocket.Conn) {
	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			if !isAcceptedError(err) {
				slog.Error("Error reading from socket", slog.String("reason", err.Error()))
			}
			return
		}

		_, err = amp.Write(data)
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
