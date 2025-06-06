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

	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			log.Fatalln("Error reading from socket:", err)
		}
		log.Printf("received: %q\n", data)

		_, err = amp.Write(data)
		if err != nil {
			log.Fatalln("Error writing to amplifier:", err)
		}

		buf := [len("-v.100\r")]byte{}
		n, err := amp.Read(buf[:])
		if err != nil {
			log.Fatalln("Error reading from amplifier:", err)
		}

		err = ws.Write(context.Background(), websocket.MessageText, buf[:n])
		if err != nil {
			log.Fatalln("Error writing back to socket:", err)
		}

		log.Printf("sent: %q\n", string(buf[:n]))
	}
}
