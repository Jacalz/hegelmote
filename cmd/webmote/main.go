package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./wasm")))
	http.Handle("/proxy", http.HandlerFunc(proxyHandler))
	http.Handle("/upnp", http.HandlerFunc(upnpHandler))

	const port = "8080"
	fmt.Println("Serving at: http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("Error when running server:", err)
	}
}
