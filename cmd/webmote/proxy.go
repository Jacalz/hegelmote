package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync/atomic"

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

	sending := atomic.Bool{}
	output := make(chan readResponse)

	go func() {
		for {
			buf := make([]byte, 32)
			n, err := amp.Read(buf)
			if sending.CompareAndSwap(true, false) {
				output <- readResponse{buf[:n], n, err}
				continue
			} else if err != nil {
				log.Fatalln("Error reading from amplifier:", err)
			}

			err = ws.Write(context.Background(), websocket.MessageText, buf[:n])
			if err != nil {
				log.Fatalln("Error writing to socket:", err)
			}
		}
	}()

	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			log.Fatalln("Error reading from socket:", err)
		}

		sending.Store(true)

		_, err = amp.Write(data)
		if err != nil {
			log.Fatalln("Error writing to amplifier:", err)
		}

		result := <-output
		if result.err != nil {
			log.Fatalln("Error reading from amplifier:", result.err)
		}

		err = ws.Write(context.Background(), websocket.MessageText, result.buf)
		if err != nil {
			log.Fatalln("Error writing back to socket:", err)
		}
	}
}

type readResponse struct {
	buf []byte
	n   int
	err error
}
