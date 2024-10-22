package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	frontUrl := os.Getenv("FRONT_URL")
	broker := NewServer()
	chatBot := NewChatBot()
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontUrl},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
	})
	handler := c.Handler(router)
	router.HandleFunc("/messages", broker.BroadcastMessage).Methods("POST")
	router.HandleFunc("/stream", broker.Stream).Methods("GET")
	router.HandleFunc("/chat", HandleConnections)
	router.HandleFunc("/chat_bot", chatBot.Chat).Methods("POST")
	go HandleMessages()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
