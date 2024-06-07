package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Broker struct to manage connected clients
type Broker struct {
	clients map[chan string]struct{}
	lock    sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan string]struct{}),
	}
}

func (b *Broker) AddClient(client chan string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.clients[client] = struct{}{}
}

func (b *Broker) RemoveClient(client chan string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	delete(b.clients, client)
	close(client)
}

func (b *Broker) Broadcast(message string) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	for client := range b.clients {
		client <- message
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan string)
	b.AddClient(clientChan)
	defer b.RemoveClient(clientChan)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Send initial message to establish the connection
	w.Write([]byte("data: Connection established\n\n"))
	flusher.Flush()

	// Listen for messages from the broker and send them to the client
	for {
		select {
		case message := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func main() {
	broker := NewBroker()

	// Register the Broker's ServeHTTP method as the handler for the /events endpoint
	http.Handle("/events", broker)

	// Additional endpoint to send messages to clients
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Query().Get("message")
		if message != "" {
			broker.Broadcast(message)
			w.Write([]byte("Message sent"))
		} else {
			http.Error(w, "Message is required", http.StatusBadRequest)
		}
	})
}
