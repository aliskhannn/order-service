package server

import "net/http"

// New creates a configured HTTP server with the given address and handler.
// This is a thin wrapper around http.Server to standardize server initialization.
func New(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
