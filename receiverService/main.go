package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var mu sync.Mutex

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println(err)
		return
	}
	redis := connect()

	cancelFuncs := make(map[string]context.CancelFunc)

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection:", err)
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		mu.Lock()
		ctx, cancel := context.WithCancel(context.Background())
		if cancelFunc, exists := cancelFuncs[userID]; exists {
			cancelFunc()
			delete(cancelFuncs, userID)
		}
		cancelFuncs[userID] = cancel
		go receiveAgent(conn, redis, userID, ctx)
		mu.Unlock()
		success(w, "Socket created successfully")

	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var requestData map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
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
		success(w, "Message sent successfully")
		fmt.Println("Good bye")
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                             // Frontend origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // Allowed HTTP methods
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // Allowed headers
	})
	handlerWithCORS := c.Handler(http.DefaultServeMux)

	fmt.Println("Server is running on http://localhost:81")
	err = http.ListenAndServe("0.0.0.0:81", handlerWithCORS)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
