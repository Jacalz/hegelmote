package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jacalz/hegelmote/internal/profile"
)

func main() {
	stop := profile.Start()
	if stop != nil {
		defer stop()
	}

	http.Handle("/", http.FileServer(http.Dir("./wasm")))
	http.Handle("/proxy", http.HandlerFunc(proxyHandler))
	http.Handle("/upnp", http.HandlerFunc(upnpHandler))

	const port = "8086"
	fmt.Println("Serving at: http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("Error when running server:", err)
	}
}
