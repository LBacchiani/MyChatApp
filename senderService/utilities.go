package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/supabase-community/supabase-go"
)

const (
	API_URL    = "https://zmyzypfirdaktluzqrkm.supabase.co"
	API_KEY    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InpteXp5cGZpcmRha3RsdXpxcmttIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzI5NjIwODIsImV4cCI6MjA0ODUzODA4Mn0.Y6f1g-xkchpwjWqV1wWCTbOaMSMc9ZNv7cbJem6NSPo"
	REDIS_URL  = "127.0.0.1"
	REDIS_PORT = "6379"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
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

func processRequest(w http.ResponseWriter, r *http.Request) Message {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	defer r.Body.Close()
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}
	return msg
}

func pushOnRedis(client *redis.Client, w http.ResponseWriter, msg Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling message:", err)
		http.Error(w, "Invalid message format", http.StatusBadRequest)
		return
	}

	err = client.Publish(ctx, msg.Receiver, msgJSON).Err()
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		http.Error(w, "Error publishing message to Redis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Status 201 Created
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Message sent successfully",
	})
}
