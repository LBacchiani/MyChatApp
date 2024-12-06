package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/supabase-community/supabase-go"
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

func connect() (*supabase.Client, *redis.Client) {
	client, err := supabase.NewClient(API_URL, API_KEY, &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", REDIS_URL, REDIS_PORT),
	})
	return client, redis
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
