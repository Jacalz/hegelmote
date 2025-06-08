package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Jacalz/hegelmote/internal/profile"
)

func main() {
	stop := profile.Start()
	if stop != nil {
		defer stop()
	}

	port := uint64(8086)
	flag.Uint64Var(&port, "port", port, "port to serve on")
	noWASM := false
	flag.BoolVar(&noWASM, "no-wasm", noWASM, "disable hosting of WASM files")
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		fmt.Printf("invalid arguments: %v\n", args)
		flag.Usage()
		return
	}

	if !noWASM {
		http.Handle("/", http.FileServer(http.Dir("./wasm")))
	}
	http.Handle("/proxy", http.HandlerFunc(proxyHandler))
	http.Handle("/upnp", http.HandlerFunc(upnpHandler))

	portString := strconv.FormatUint(port, 10)

	fmt.Printf("Serving at: http://localhost:%s\n", portString)

	server := http.Server{Addr: ":" + portString, ReadTimeout: time.Second, WriteTimeout: time.Second}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error when running server:", err)
	}
}
