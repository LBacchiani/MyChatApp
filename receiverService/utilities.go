package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/go-redis/redis/v8"
)

const (
	API_URL    = "https://zmyzypfirdaktluzqrkm.supabase.co"
	API_KEY    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InpteXp5cGZpcmRha3RsdXpxcmttIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzI5NjIwODIsImV4cCI6MjA0ODUzODA4Mn0.Y6f1g-xkchpwjWqV1wWCTbOaMSMc9ZNv7cbJem6NSPo"
	REDIS_URL  = "127.0.0.1"
	REDIS_PORT = "6379"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins (for simplicity, modify for production)
		return true
	},
}

func connect() *redis.Client {

	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT")),
	})
	return redis
}

func receiveAgent(conn *websocket.Conn, client *redis.Client, id string, ctx context.Context) {
	fmt.Println("RECEIVE AGENT STARTED")
	channel := client.Subscribe(ctx, id)
	for {
		msg, err := channel.ReceiveMessage(ctx)
		if err != nil {
			conn.Close()
			return
		}
		fmt.Println("Received message: %s\n" + msg.Payload)
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("Error writing JSON response:", err)
			conn.Close()
			return
		}
	}
}

func success(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Status 201 Created
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": message,
	})
}
