//go:build bundled

package main

import (
	"embed"
	"net/http"
)

//go:embed wasm/*
var wasm embed.FS

// Serving WASM files embedded in the binary.
func serveWASM() {
	http.Handle("/", http.FileServer(http.FS(wasm)))
}
