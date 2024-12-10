package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/supabase-community/supabase-go"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
	IsRead   bool   `json:"isRead"`
}

type AckMessage struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Type     string `json:"type"`
}

func connect() (*supabase.Client, *redis.Client) {
	client, err := supabase.NewClient(os.Getenv("API_URL"), os.Getenv("API_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT")),
	})
	return client, redis
}

func processRequest(w http.ResponseWriter, r *http.Request) interface{} {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	defer r.Body.Close()
	var msg Message
	if err := json.Unmarshal(body, &msg); err == nil && msg.Content != "" {
		return msg
	}
	var ackMsg AckMessage
	if err := json.Unmarshal(body, &ackMsg); err == nil {
		return ackMsg
	}

	http.Error(w, "Invalid request body format", http.StatusBadRequest)
	return nil
}

func pushOnRedis(client *redis.Client, w http.ResponseWriter, msg interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var receiver string
	var msgJSON []byte
	var err error
	switch m := msg.(type) {
	case Message:
		receiver = m.Receiver
		msgJSON, err = json.Marshal(m)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			http.Error(w, "Invalid message format", http.StatusBadRequest)
			return
		}
	case AckMessage:
		receiver = m.Receiver
		msgJSON, err = json.Marshal(m)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			http.Error(w, "Invalid message format", http.StatusBadRequest)
			return
		}
	}
	err = client.Publish(ctx, receiver, msgJSON).Err()
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		http.Error(w, "Error publishing message to Redis", http.StatusInternalServerError)
		return
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
