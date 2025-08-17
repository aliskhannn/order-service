package handlers

import (
	"encoding/json"
	"net/http"
)

// WriteError writes an error response in JSON format with the given status code.
// The response body has the structure: {"error": "<msg>"}.
func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// WriteJSON writes a JSON response with the given status code and arbitrary data.
// The 'data' argument can be any Go value that is JSON serializable.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
