package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Jacalz/hegelmote/internal/profile"
)

func main() {
	stop := profile.Start()
	if stop != nil {
		defer stop()
	}

	http.Handle("/", http.FileServer(http.Dir("../../wasm")))
	http.Handle("/proxy", http.HandlerFunc(proxyHandler))
	http.Handle("/upnp", http.HandlerFunc(upnpHandler))

	port := uint64(8086)
	flag.Uint64Var(&port, "port", port, "port to serve on")
	flag.Parse()
	portString := strconv.FormatUint(port, 10)

	fmt.Printf("Serving at: http://localhost:%s\n", portString)
	err := http.ListenAndServe(":"+portString, nil)
	if err != nil {
		log.Fatalln("Error when running server:", err)
	}
}
