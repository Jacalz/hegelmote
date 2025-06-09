//go:build !bundled

package main

import "net/http"

// Serving WASM files from ./wasm directory next to binary.
func serveWASM() {
	http.Handle("/", http.FileServer(http.Dir("./wasm")))
}
