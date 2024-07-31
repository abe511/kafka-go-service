package main

import (
	"log"
	"net/http"
	"kafka-go-service/handlers"
	"kafka-go-service/kafkaservice"
)


func main() {

	kafkaservice.InitKafka()

	kafkaservice.RunConsumer()

	// start a server with two endpoints
	router := http.NewServeMux()
	router.HandleFunc("GET /stats", handlers.GetStats)
	router.HandleFunc("POST /message", handlers.ReceiveMessage)

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	log.Printf("Server started on port%v\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}