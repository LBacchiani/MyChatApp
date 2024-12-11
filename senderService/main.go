package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println(err)
		return
	}
	supabase, redis := connect()
	// messageQueue :=

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := processRequest(w, r)
		push := false
		switch m := (msg).(type) {
		case Message:
			fmt.Println("Message:", m)
			_, _, err := supabase.From("Message").Insert(m, false, "", "", "1").Execute()
			if err != nil {
				http.Error(w, "Error sending message: "+err.Error(), http.StatusInternalServerError)
			}
			push = m.Receiver != m.Sender

		case AckMessage:
			fmt.Println("Ack:", m)
			_, _, err := supabase.From("Message").Update(map[string]interface{}{"isRead": true, "sender": m.Receiver}, "", "").Eq("isRead", "false").Eq("sender", m.Receiver).Execute()
			push = true
			if err != nil {
				http.Error(w, "Error updating messages: "+err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			fmt.Println("Invalid message")
			http.Error(w, "Invalid message format", http.StatusBadRequest)
			return
		}
		if push {
			pushOnRedis(redis, w, msg)
		}
		success(w, "Message sent successfully")
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                             // Frontend origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // Allowed HTTP methods
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // Allowed headers
	})
	handlerWithCORS := c.Handler(http.DefaultServeMux)

	fmt.Println("Server is running on http://localhost:80")
	err = http.ListenAndServe("0.0.0.0:80", handlerWithCORS)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
