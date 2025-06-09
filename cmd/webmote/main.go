package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Jacalz/hegelmote/internal/profile"
)

func main() {
	stop := profile.Start()
	if stop != nil {
		defer stop()
	}

	portNumber := uint64(8086)
	flag.Uint64Var(&portNumber, "port", portNumber, "port to serve on")
	noWASM := false
	flag.BoolVar(&noWASM, "no-wasm", noWASM, "disable hosting of WASM files")
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		fmt.Printf("invalid arguments: %v\n", args)
		flag.Usage()
		return
	}

	logfile, err := os.CreateTemp(".", "hegelmote-*.log")
	if err != nil {
		log.Fatalln("Error creating log file:", err)
	}
	defer logfile.Close()

	logger := slog.NewTextHandler(io.MultiWriter(os.Stdout, logfile), nil)
	slog.SetDefault(slog.New(logger))

	if !noWASM {
		http.Handle("/", http.FileServer(http.Dir("./wasm")))
	}
	http.Handle("/proxy", http.HandlerFunc(proxyHandler))
	http.Handle("/upnp", http.HandlerFunc(upnpHandler))

	port := strconv.FormatUint(portNumber, 10)
	fmt.Printf("Serving at: http://localhost:%s\n", port)

	const timeout = time.Second
	server := http.Server{Addr: ":" + port, ReadTimeout: timeout, WriteTimeout: timeout, ErrorLog: slog.NewLogLogger(logger, slog.LevelInfo)}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error when running server:", err)
	}
}
