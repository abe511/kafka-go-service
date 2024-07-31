package handlers

import (
	"encoding/json"
	"kafka-go-service/models"
	"kafka-go-service/kafkaservice"
	"net/http"
)


func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("message received\n"))
	var msg models.Message

	// decode the msg from the request body json
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// forward the msg to kafka producer
	err = kafkaservice.SendToKafka(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// respond with msg and status 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}



func GetStats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " stats\n"))
}
