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
	fmt.Println("RECEIVE AGENT FOR " + id + " STARTED")
	channel := client.Subscribe(ctx, id)
	defer func() {
		_ = channel.Close()
		_ = conn.Close()
	}()
	for {
		msg, err := channel.ReceiveMessage(ctx)
		if err != nil {
			conn.Close()
			return
		}
		fmt.Println("Received message\n" + msg.Payload)
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
