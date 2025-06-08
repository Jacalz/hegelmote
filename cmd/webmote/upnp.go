package main

import (
	"context"
	"log"
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
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Failed to accept upnp socket:", err)
		return
	}
	defer ws.CloseNow()

	devices, err := upnp.LookUpDevices()
	err = wsjson.Write(context.Background(), ws, upnpResponse{devices, err})
	if err != nil {
		log.Println("Failed to write upnp devices:", err)
	}
}
