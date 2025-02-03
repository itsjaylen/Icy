package admin

import (
	"encoding/json"
	"net/http"
)

// GetStatusHandler handles the /admin/status route
func GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
