package admin

import (
	"IcyAPI/internal/utils"
	"net/http"
)

// GetStatusHandler handles the /admin/status route
func GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}
