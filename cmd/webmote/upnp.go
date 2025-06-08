package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Jacalz/hegelmote/internal/upnp"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type upnpResponse struct {
	Devices []upnp.DiscoveredDevice `json:"devices"`
	Err     error                   `json:"error"`
}

func upnpHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("New UPNP lookup", slog.String("remote", r.RemoteAddr))

	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("Failed to accept upnp socket:", slog.String("error", err.Error()))
		return
	}
	defer ws.CloseNow()

	devices, err := upnp.LookUpDevices()
	err = wsjson.Write(context.Background(), ws, upnpResponse{devices, err})
	if err != nil {
		slog.Error("Failed to write upnp devices:", slog.String("error", err.Error()))
	}
}
