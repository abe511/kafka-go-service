package main

import (
	"log"
	"net/http"
	"github.com/abe511/kafka-service/handlers"
)


func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /stats", handlers.GetStats)
	router.HandleFunc("POST /message", handlers.SendMessage)

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	log.Printf("Server started on port%v\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}