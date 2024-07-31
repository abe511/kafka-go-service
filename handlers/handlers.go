package handlers

import (
	"encoding/json"
	"kafka-go-service/models"
	"kafka-go-service/kafkaservice"
	"kafka-go-service/database"
	"net/http"
)


func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message

	// extract the message from the request body
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// store the message in the db
	err = database.StoreMessage(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// forward the message to kafka producer
	err = kafkaservice.SendToKafka(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// respond with the message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// get current statistics on sent and processed messages total
func GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := database.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
