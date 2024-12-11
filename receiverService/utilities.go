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
	initStream(client, id, ctx)
	defer conn.Close()
	lastMessageID := "0"
	for {
		streams, err := client.XRead(ctx, &redis.XReadArgs{Streams: []string{id, lastMessageID}, Count: 1, Block: 0}).Result()

		if err != nil {
			fmt.Println("Error reading from Redis Stream:", err)
			return
		}
		for _, stream := range streams {
			for _, message := range stream.Messages {
				for _, value := range message.Values {
					fmt.Println("Received message\n", value)
					lastMessageID = message.ID
					_, err = client.XDel(ctx, id, message.ID).Result()
					if err != nil {
						fmt.Println("Error removing from Redis Stream:", err)
						return
					}
					if value != id {
						if err := conn.WriteJSON(value); err != nil {
							fmt.Println("Error writing JSON response:", err)
							return
						}
					}
				}
			}
		}
	}
}

func initStream(client *redis.Client, id string, ctx context.Context) {
	_, err := client.XAdd(ctx, &redis.XAddArgs{Stream: id, Values: map[string]interface{}{"key": id}}).Result()
	if err != nil {
		fmt.Println("Error writing on stream id")
	}
}

func deleteStreamOnCancel(client *redis.Client, streamID string, ctx context.Context) {
	<-ctx.Done()
	fmt.Println("Context canceled. Deleting stream:", streamID)
	_, err := client.Del(ctx, streamID).Result()
	if err != nil {
		fmt.Println("Error deleting stream:", err)
	} else {
		fmt.Println("Stream deleted successfully.")
	}
}

func success(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": message,
	})
}
