package main

import (
	"log"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	kafkaBrokers = "localhost:9092"
	kafkaTopic   = "my_topic"
)

func main() {
	// Initializing a connection with Kafka
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokers})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	topic := kafkaTopic
	// Sending 2000 messages to Kafka
	for i := 1; i <= 2000; i++ {
		message := []byte("Message " + strconv.Itoa(i))
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: message,
		}, nil)

		time.Sleep(time.Millisecond)
	}

	// Waiting for the completion of sending messages
	producer.Flush(15 * 1000)
}
