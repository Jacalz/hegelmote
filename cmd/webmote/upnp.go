package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Jacalz/hegelmote/internal/upnp"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func upnpHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatalln("Failed to accept upnp socket:", err)
	}
	defer ws.CloseNow()

	devices, err := upnp.LookUpDevices()
	if err != nil {
		log.Fatalln("Failed to look up UPnP devices:", err)
	}

	err = wsjson.Write(context.Background(), ws, devices)
	if err != nil {
		log.Fatalln("Failed to write upnp devices:", err)
	}
}
