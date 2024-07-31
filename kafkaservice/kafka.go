package kafkaservice

import (
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafka-go-service/models"
)

var Producer *kafka.Producer
var Consumer *kafka.Consumer

func InitKafka() {
	var err error

	// initialize a producer
	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
	})

	if err != nil {
		log.Fatal(err)
	}

	// initialize a consumer
	Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
		"group.id": "test_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}
}

// produce a message and add it to the queue
func SendToKafka(msg *models.Message) error {
	topic := "test_topic"
	msgBytes, _ := json.Marshal(msg)

	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: msgBytes,
	}, nil)
	if err != nil {
		return err
	}

	return err
}


// subscribe to a topic and run the consumer in a separate goroutine
func RunConsumer() {
	topic := "test_topic"
	err := Consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal(err)
	}

	go consumeMessages()
}

// run the consumer in an infinite loop, trying to consume the latest message
func consumeMessages() {
	for {
		msg, err := Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v. msg: %v\n", err, msg)
			continue
		}
		
		var receivedMsg models.Message

		err = json.Unmarshal(msg.Value, &receivedMsg)
		if err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		log.Printf("Received message: %+v\n", receivedMsg)
	}

}