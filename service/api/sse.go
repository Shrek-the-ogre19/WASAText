package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	clients   = make(map[chan []byte]bool)
	clientsMu sync.RWMutex
)

func BroadcastRefresh() {
	msg := []byte(`{"action": "refresh"}`)

	clientsMu.RLock()
	defer clientsMu.RUnlock()

	for clientChan := range clients {
		select {
		case clientChan <- msg:
		default:
		}
	}
}

func SSEHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("📡 New SSE client connecting...")

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set SSE required headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	// Check if flusher is available
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Create a UNIQUE channel for THIS client (using chan string)
	clientChan := make(chan []byte, 10)

	// Register this client in the global map
	clientsMu.Lock()
	clients[clientChan] = true
	clientsMu.Unlock()

	// Cleanup when client disconnects
	defer func() {
		clientsMu.Lock()
		delete(clients, clientChan)
		clientsMu.Unlock()
		close(clientChan)
		fmt.Println("👋 Client disconnected. Total clients:", len(clients))
	}()

	// Send initial connection success message
	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: {\"status\":\"connected\"}\n\n")
	flusher.Flush()
	fmt.Println("✅ Initial message sent. Total clients:", len(clients))

	// Send keep-alive ping every 15 seconds to prevent timeout
	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	// Main event loop - keeps connection open for THIS client
	for {
		select {
		case msg := <-clientChan:
			// Send broadcast message to THIS client
			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
			fmt.Println("📤 Message sent to client")

		case <-pingTicker.C:
			// Send a comment (ping) to keep connection alive
			fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()
			fmt.Println("💓 Keep-alive ping sent")

		case <-r.Context().Done():
			// Client disconnected
			fmt.Println("❌ Client context cancelled")
			return
		}
	}
}
