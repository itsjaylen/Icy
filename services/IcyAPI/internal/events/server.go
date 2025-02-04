package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	logger "itsjaylen/IcyLogger"
)

// AdminEvent represents an event triggered by admin actions
type AdminEvent struct {
	Type    string `json:"type"`
	Action  string `json:"action"`
	Target  string `json:"target"`  // The user or object the action applies to
	Message string `json:"message"` // Any additional message or reason for the action
}

// PublishAdminEvent publishes an admin event to all clients
func (s *EventServer) PublishAdminEvent(action, target, message string) {
	event := AdminEvent{
		Type:    "admin",
		Action:  action,
		Target:  target,
		Message: message,
	}
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error.Printf("Error marshalling admin event: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		if client.subscribe["admin"] || client.subscribe["*"] {
			select {
			case client.send <- string(data):
			default:
				logger.Error.Printf("Dropping message for slow client %d", client.id)
			}
		}
	}
}

// EventServer manages SSE clients and event broadcasting
type EventServer struct {
	mu      sync.Mutex
	clients map[int]*Client
	counter int
	Host    string
	Port    string
}

// Event represents an SSE event payload
type Event struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Action  string `json:"action"`
	Content string `json:"content,omitempty"`
}

// Publishes events to all subscribed clients
func (s *EventServer) Publish(eventType, name, action string) {
	event := Event{Type: eventType, Name: name, Action: action}
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error.Printf("Error marshalling event: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		if client.subscribe["*"] || client.subscribe[eventType] {
			select {
			case client.send <- string(data):
			default:
				logger.Error.Printf("Dropping message for slow client %d", client.id)
			}
		}
	}
}

// NewEventServer initializes a new SSE server
func NewEventServer(host, port string) *EventServer {
	return &EventServer{
		clients: make(map[int]*Client),
		Host:    host,
		Port:    port,
	}
}

// Start begins listening for connections and requests
func (s *EventServer) Start() error {
	http.HandleFunc("/events", s.addClient)
	http.HandleFunc("/add_user", s.addUserHandler)
	http.HandleFunc("/admin_event", s.adminEventHandler) // Add new route here

	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	logger.Info.Printf("Event server starting on %s...", address)
	return http.ListenAndServe(address, nil)
}

// Shutdown gracefully closes all client connections
func (s *EventServer) Shutdown() {
	logger.Info.Println("Shutting down event server...")
	s.mu.Lock()
	for _, client := range s.clients {
		close(client.send)
	}
	s.clients = make(map[int]*Client)
	s.mu.Unlock()
}
