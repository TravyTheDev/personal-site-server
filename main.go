package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	frontUrl := os.Getenv("FRONT_URL")
	port := os.Getenv("PORT")
	broker := NewServer()
	chatBot := NewChatBot()
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/my_site_api").Subrouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontUrl},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
	})
	handler := c.Handler(router)
	subRouter.HandleFunc("/messages", broker.BroadcastMessage).Methods("POST")
	subRouter.HandleFunc("/stream", broker.Stream).Methods("GET")
	subRouter.HandleFunc("/chat", HandleConnections)
	subRouter.HandleFunc("/chat_bot", chatBot.Chat).Methods("POST")
	go HandleMessages()

	log.Println("Starting server on", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
