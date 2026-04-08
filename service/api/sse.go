package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var clients = make(map[chan string]bool)
var register = make(chan chan string)
var unregister = make(chan chan string)
var broadcast = make(chan string)

func broadcaster() {
	log.Println("Broadcaster started")
	for {
		select {
		case client := <-register:
			log.Printf("Registering new client. Total before: %d", len(clients))
			clients[client] = true
			log.Printf("Total clients after: %d", len(clients))

		case client := <-unregister:
			log.Printf("Unregistering client. Total before: %d", len(clients))
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client) // Close the channel to signal the handler
				log.Printf("Client removed. Total after: %d", len(clients))
			}

		case message := <-broadcast:
			log.Printf("Broadcasting message: %q to %d clients", message, len(clients))
			for client := range clients {
				// ALWAYS use a select with a small timeout, never default
				select {
				case client <- message:
					log.Printf("Message sent to client")
				case <-time.After(1 * time.Second):
					// Client is slow, but don't remove it immediately
					log.Printf("Client slow, skipping this message")
				}
			}
		}
	}
}

func (rt *_router) sseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create buffered channel
	messageChan := make(chan string, 10)

	// Register the client FIRST
	register <- messageChan
	log.Println("SSE: Client registered")

	// Wait a moment for registration to complete
	time.Sleep(10 * time.Millisecond)

	// Send initial comment
	if _, err := fmt.Fprintf(w, ": connected\n\n"); err != nil {
		log.Printf("SSE: Failed to write initial data: %v", err)
		return
	}
	flusher.Flush()
	log.Println("SSE: Initial data sent")

	// Cleanup
	defer func() {
		log.Println("SSE: Defer unregister")
		unregister <- messageChan
	}()

	// Main loop
	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				log.Println("SSE: Channel closed")
				return
			}
			log.Printf("SSE: Sending message: %s", message)
			if _, err := fmt.Fprintf(w, "data: %s\n\n", message); err != nil {
				log.Printf("SSE: Error writing: %v", err)
				return
			}
			flusher.Flush()
			log.Println("SSE: Message sent")

		case <-r.Context().Done():
			log.Printf("SSE: Context done, client disconnected: %v", r.Context().Err())
			return
		}
	}
}

func (rt *_router) sendMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message := r.URL.Query().Get("message")
	if message == "" {
		message = "Hello from server!"
	}

	// Broadcast to all clients
	broadcast <- message

	w.Write([]byte("Message sent to all clients: " + message))
}
