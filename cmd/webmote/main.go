package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/coder/websocket"
)

func proxy(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatalln("Failed to accept websocket:", err)
	}
	defer ws.CloseNow()

	amp, err := net.Dial("tcp", "192.168.1.251:50001")
	if err != nil {
		log.Fatalln("Failed to connect to amplifier:", err)
	}
	defer amp.Close()

	for {
		_, data, err := ws.Read(context.Background())
		if err != nil {
			log.Fatalln("Error reading from websocket:", err)
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
			log.Fatalln("Error writing back to websocket:", err)
		}

		log.Printf("sent: %q\n", string(buf[:n]))
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./wasm")))
	http.Handle("/proxy", http.HandlerFunc(proxy))

	const port = "8080"
	fmt.Println("Serving at: http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("Error when running server:", err)
	}
}
