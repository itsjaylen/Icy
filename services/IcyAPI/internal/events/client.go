package events

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Client represents an SSE subscriber
type Client struct {
	id        int
	send      chan string
	subscribe map[string]bool 
}

// Handles new client connections
func (s *EventServer) addClient(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	events := strings.Split(r.URL.Query().Get("events"), ",")
	subscribe := make(map[string]bool)
	for _, event := range events {
		trimmedEvent := strings.TrimSpace(event)
		if trimmedEvent == "*" {
			subscribe["*"] = true
			break
		}
		subscribe[trimmedEvent] = true
	}

	s.mu.Lock()
	s.counter++
	client := &Client{id: s.counter, send: make(chan string, 10), subscribe: subscribe}
	s.clients[client.id] = client
	s.mu.Unlock()

	log.Printf("Client %d connected", client.id)

	defer func() {
		s.mu.Lock()
		delete(s.clients, client.id)
		s.mu.Unlock()
		close(client.send)
		log.Printf("Client %d disconnected", client.id)
	}()

	for msg := range client.send {
		_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
		if err != nil {
			return
		}
		flusher.Flush()
	}
}
