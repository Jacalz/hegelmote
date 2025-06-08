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
		log.Println("Failed to accept proxy socket:", err)
		return
	}
	defer ws.CloseNow()

	amp, err := connect(ws)
	if err != nil {
		log.Println("Failed to connect to amplifier:", err)
		return
	}
	defer amp.Close()

	go forwardFromAmplifier(amp, ws)
	forwardFromClient(amp, ws)

	log.Println("Shutting down connection to client and amplifier")
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
			log.Println("Error reading from amplifier:", err)
			return
		}

		err = ws.Write(context.Background(), websocket.MessageText, buf[:n])
		if err != nil {
			log.Println("Error writing to socket:", err)
			return
		}
	}
}

func forwardFromClient(amp net.Conn, ws *websocket.Conn) {
	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			log.Println("Error reading from socket:", err)
			return
		}

		_, err = amp.Write(data)
		if err != nil {
			log.Println("Error writing to amplifier:", err)
			return
		}
	}
}
