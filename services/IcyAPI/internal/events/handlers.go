package events

import (
	logger "itsjaylen/IcyLogger"
	"net/http"
)

func (s *EventServer) addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}
	logger.Info.Printf("User added: %s", username)
	s.Publish("user", username, "added")
	w.WriteHeader(http.StatusOK)
}
