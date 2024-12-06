package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/cors"
)

var mu sync.Mutex

func main() {
	_, redis := connect()

	cancelFuncs := make(map[string]context.CancelFunc)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},         // Frontend origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // Allowed HTTP methods
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // Allowed headers
	})

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		// First, upgrade the HTTP request to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection:", err)
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		// Extract the user_id from the query parameters
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		// Create a context and cancel function for the user connection
		ctx, cancel := context.WithCancel(context.Background())

		mu.Lock()
		cancelFuncs[userID] = cancel
		mu.Unlock()

		// Handle the WebSocket connection in a separate goroutine
		go receiveAgent(conn, redis, userID, ctx)

		// Respond back with a success message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Status 201 Created
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "WebSocket connection established successfully",
		})
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var requestData map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		// Access the user_id from the map
		userID, ok := requestData["user_id"].(string)
		if !ok {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		mu.Lock()
		if cancelFunc, exists := cancelFuncs[userID]; exists {
			cancelFunc()
			delete(cancelFuncs, userID)
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Status 201 Created
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Message sent successfully",
		})
		fmt.Println("Good bye")
	})

	handlerWithCORS := c.Handler(http.DefaultServeMux)

	fmt.Println("Server is running on http://localhost:81")
	err := http.ListenAndServe(":81", handlerWithCORS)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
