package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse is a helper function that encodes a response to JSON and writes it to the
// ResponseWriter.
// It handles errors and sends an appropriate HTTP status code if encoding fails.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
