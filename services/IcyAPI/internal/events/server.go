package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	logger "itsjaylen/IcyLogger"
)

// AdminEvent represents an event triggered by admin actions.
type AdminEvent struct {
	Type    string `json:"type"`
	Action  string `json:"action"`
	Target  string `json:"target"`
	Message string `json:"message"`
}

// PublishAdminEvent publishes an admin event to all clients.
func (server *EventServer) PublishAdminEvent(action, target, message string) {
	event := AdminEvent{
		Type:    "admin",
		Action:  action,
		Target:  target,
		Message: message,
	}
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error.Printf("Error marshaling admin event: %v", err)

		return
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	for _, client := range server.clients {
		if client.subscribe["admin"] || client.subscribe["*"] {
			select {
			case client.send <- string(data):
			default:
				logger.Error.Printf("Dropping message for slow client %d", client.id)
			}
		}
	}
}

// EventServer manages SSE clients and event broadcasting.
type EventServer struct {
	clients map[int]*Client
	Host    string
	Port    string
	counter int
	mu      sync.Mutex
}

// Event represents an SSE event payload.
type Event struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Action  string `json:"action"`
	Content string `json:"content,omitempty"`
}

// Publishes events to all subscribed clients.
func (server *EventServer) Publish(eventType, name, action string) {
	event := Event{Type: eventType, Name: name, Action: action}
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error.Printf("Error marshaling event: %v", err)

		return
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	for _, client := range server.clients {
		if client.subscribe["*"] || client.subscribe[eventType] {
			select {
			case client.send <- string(data):
			default:
				logger.Error.Printf("Dropping message for slow client %d", client.id)
			}
		}
	}
}

// NewEventServer initializes a new SSE server.
func NewEventServer(host, port string) *EventServer {
	return &EventServer{
		clients: make(map[int]*Client),
		Host:    host,
		Port:    port,
	}
}

// Start begins listening for connections and requests.
func (server *EventServer) Start() error {
	http.HandleFunc("/events", server.addClient)
	http.HandleFunc("/add_user", server.addUserHandler)
	http.HandleFunc("/admin_event", server.adminEventHandler) // Add new route here

	address := fmt.Sprintf("%s:%s", server.Host, server.Port)
	logger.Info.Printf("Event server starting on %s...", address)

	return http.ListenAndServe(address, nil)
}

// Shutdown gracefully closes all client connections.
func (server *EventServer) Shutdown() error {
	logger.Info.Println("Shutting down event server...")
	server.mu.Lock()
	for _, client := range server.clients {
		close(client.send)
	}
	server.clients = make(map[int]*Client)
	server.mu.Unlock()

	return nil
}
