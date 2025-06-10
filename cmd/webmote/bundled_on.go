//go:build bundled

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed wasm/*
var wasm embed.FS

// Serving WASM files embedded in the binary.
func serveWASM() {
	files, _ := fs.Sub(wasm, "wasm")
	http.Handle("/", http.FileServer(http.FS(files)))
}
