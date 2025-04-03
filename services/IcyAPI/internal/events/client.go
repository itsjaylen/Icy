package events

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Client represents an SSE subscriber. TODO: Redo event server.
type Client struct {
	send      chan string
	subscribe map[string]bool
	id        int
}

// Handles new client connections.
func (server *EventServer) addClient(writer http.ResponseWriter, request *http.Request) {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "Streaming unsupported", http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	events := strings.Split(request.URL.Query().Get("events"), ",")
	subscribe := make(map[string]bool)
	for _, event := range events {
		trimmedEvent := strings.TrimSpace(event)
		if trimmedEvent == "*" {
			subscribe["*"] = true

			break
		}
		subscribe[trimmedEvent] = true
	}

	server.mu.Lock()
	server.counter++
	client := &Client{id: server.counter, send: make(chan string, 10), subscribe: subscribe}
	server.clients[client.id] = client
	server.mu.Unlock()

	log.Printf("Client %d connected", client.id)

	defer func() {
		server.mu.Lock()
		delete(server.clients, client.id)
		server.mu.Unlock()
		close(client.send)
		log.Printf("Client %d disconnected", client.id)
	}()

	for msg := range client.send {
		_, err := fmt.Fprintf(writer, "data: %s\n\n", msg)
		if err != nil {
			return
		}
		flusher.Flush()
	}
}
