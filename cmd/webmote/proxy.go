package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/coder/websocket"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("New proxy connection", slog.String("source", r.RemoteAddr))

	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("Failed to accept proxy socket:", slog.String("error", err.Error()))
		return
	}
	defer ws.CloseNow()

	amp, err := connect(ws)
	if err != nil {
		slog.Error("Failed to connect to amplifier:", slog.String("error", err.Error()))
		return
	}
	defer amp.Close()

	go forwardFromAmplifier(amp, ws)
	forwardFromClient(amp, ws)

	slog.Info("Shutting down connection to client and amplifier")
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
			slog.Error("Error reading from amplifier:", slog.String("error", err.Error()))
			return
		}

		err = ws.Write(context.Background(), websocket.MessageText, buf[:n])
		if err != nil {
			slog.Error("Error writing to socket:", slog.String("error", err.Error()))
			return
		}
	}
}

func forwardFromClient(amp net.Conn, ws *websocket.Conn) {
	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			slog.Error("Error reading from socket:", slog.String("error", err.Error()))
			return
		}

		_, err = amp.Write(data)
		if err != nil {
			slog.Error("Error writing to amplifier:", slog.String("error", err.Error()))
			return
		}
	}
}
