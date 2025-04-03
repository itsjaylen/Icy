package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse is a helper function that encodes a response to JSON and writes it to the
// ResponseWriter.
// It handles errors and sends an appropriate HTTP status code if encoding fails.
func WriteJSONResponse(writer http.ResponseWriter, statusCode int, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
	}
}
