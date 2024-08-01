package kafkaservice

import (
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafka-go-service/models"
	"kafka-go-service/database"
)

var Producer *kafka.Producer
var Consumer *kafka.Consumer
var KafkaTopic string = os.Getenv("KAFKA_TOPIC")

func InitKafka() {

	var configProducer kafka.ConfigMap
	var configConsumer kafka.ConfigMap

	brokers := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	if os.Getenv("ENV") == "global" {
		caPem := os.Getenv("KAFKA_CA_PEM")
		serviceCert := os.Getenv("KAFKA_SERVICE_CERT")
		serviceKey := os.Getenv("KAFKA_SERVICE_KEY")
		
		configProducer = kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"security.protocol": "SSL",
			"ssl.ca.location": caPem,
			"ssl.certificate.location": serviceCert,
			"ssl.key.location": serviceKey,
		}
		configConsumer = kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"security.protocol": "SSL",
			"ssl.ca.location": caPem,
			"ssl.certificate.location": serviceCert,
			"ssl.key.location": serviceKey,
			"group.id": groupID,
			"auto.offset.reset": "earliest",
		}
	} else {
		configProducer = kafka.ConfigMap{
				"bootstrap.servers": brokers,
			}
		configConsumer = kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"group.id": "test_group",
			"auto.offset.reset": "earliest",
		}
	}

	var err error

	// initialize a producer
	Producer, err = kafka.NewProducer(&configProducer)
	if err != nil {
		log.Fatal(err)
	}

	// initialize a consumer
	Consumer, err = kafka.NewConsumer(&configConsumer)
	if err != nil {
		log.Fatal(err)
	}
}

// produce a message and add it to the queue
func SendToKafka(msg *models.Message) error {
	msgBytes, _ := json.Marshal(msg)

	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &KafkaTopic, Partition: kafka.PartitionAny},
		Value: msgBytes,
	}, nil)
	if err != nil {
		return err
	}

	return err
}


// subscribe to a topic and run the consumer in a separate goroutine
func RunConsumer() {
	err := Consumer.SubscribeTopics([]string{KafkaTopic}, nil)
	if err != nil {
		log.Fatal(err)
	}

	go consumeMessages()
}

// run the consumer in an infinite loop
// consume the latest message and update its processed status
func consumeMessages() {
	for {
		msg, err := Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v. msg: %v\n", err, msg)
			continue
		}
		
		var receivedMsg models.Message
		// parse json
		err = json.Unmarshal(msg.Value, &receivedMsg)
		if err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		log.Printf("Message %v received: %v\n", receivedMsg.ID, receivedMsg.Content)
		// mark the messages as processed
		err = database.UpdateMessageStatus(receivedMsg.ID)
		if err != nil {
			log.Printf("Error updating message status: %v\n", err)
			continue
		}

		log.Printf("Message %v status: processed\n", receivedMsg.ID)
	}
}