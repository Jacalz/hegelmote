package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/coder/websocket"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatalln("Failed to accept proxy socket:", err)
	}
	defer ws.CloseNow()

	_, host, err := ws.Read(context.Background())
	if err != nil {
		log.Fatalln("Error reading host from socket:", err)
	}

	amp, err := net.Dial("tcp", string(host)+":50001")
	if err != nil {
		log.Fatalln("Failed to connect to amplifier:", err)
	}
	defer amp.Close()

	go forwardFromAmplifier(amp, ws)
	forwardFromClient(amp, ws)
}

func forwardFromAmplifier(amp net.Conn, ws *websocket.Conn) {
	buf := make([]byte, 32)
	for {
		n, err := amp.Read(buf)
		if err != nil {
			log.Fatalln("Error reading from amplifier:", err)
		}

		err = ws.Write(context.Background(), websocket.MessageText, buf[:n])
		if err != nil {
			log.Fatalln("Error writing to socket:", err)
		}
	}
}

func forwardFromClient(amp net.Conn, ws *websocket.Conn) {
	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			log.Fatalln("Error reading from socket:", err)
		}

		_, err = amp.Write(data)
		if err != nil {
			log.Fatalln("Error writing to amplifier:", err)
		}
	}
}
