package events

import (
	"net/http"

	logger "itsjaylen/IcyLogger"
)

// Placeholder handler for user events.
func (server *EventServer) addUserHandler(writer http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("username")
	if username == "" {
		http.Error(writer, "Missing username", http.StatusBadRequest)

		return
	}
	logger.Info.Printf("User added: %s", username)
	server.Publish("user", username, "added")
	writer.WriteHeader(http.StatusOK)
}

// Placeholder handler for admin events.
func (server *EventServer) adminEventHandler(writer http.ResponseWriter, request *http.Request) {
	action := request.URL.Query().Get("action")
	target := request.URL.Query().Get("target")
	message := request.URL.Query().Get("message")

	if action == "" || target == "" {
		http.Error(writer, "Missing action or target", http.StatusBadRequest)

		return
	}

	logger.Info.Printf("Admin event triggered: %s on %s", action, target)
	server.PublishAdminEvent(action, target, message)
	writer.WriteHeader(http.StatusOK)
}
