package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	supabase, redis := connect()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},         // Frontend origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // Allowed HTTP methods
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // Allowed headers
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := processRequest(w, r)
		fmt.Println("Message to be inserted:", msg)
		_, _, err := supabase.From("Message").Insert(msg, false, "", "", "1").Execute()
		if err != nil {
			http.Error(w, "Error sending message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if msg.Receiver != msg.Sender {
			pushOnRedis(redis, w, msg)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Status 201 Created
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Message sent successfully",
		})
	})

	handlerWithCORS := c.Handler(http.DefaultServeMux)

	fmt.Println("Server is running on http://localhost:80")
	err := http.ListenAndServe(":80", handlerWithCORS)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
