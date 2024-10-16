package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	broker := NewServer()
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4321"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
	})
	handler := c.Handler(router)
	router.HandleFunc("/messages", broker.BroadcastMessage).Methods("POST")
	router.HandleFunc("/stream", broker.Stream).Methods("GET")
	router.HandleFunc("/chat", HandleConnections)
	go HandleMessages()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
