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

// Placeholder handler for admin events
func (s *EventServer) adminEventHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	target := r.URL.Query().Get("target")
	message := r.URL.Query().Get("message")

	if action == "" || target == "" {
		http.Error(w, "Missing action or target", http.StatusBadRequest)
		return
	}

	logger.Info.Printf("Admin event triggered: %s on %s", action, target)
	s.PublishAdminEvent(action, target, message)
	w.WriteHeader(http.StatusOK)
}
