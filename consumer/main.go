package main

import (
	"fmt"
	"logs-mgmt/db"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

var topic = db.KAFKA_TOPIC

const consumerGroupId = "logs-consumer-group"

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          consumerGroupId,
		"auto.offset.reset": "smallest",
	}
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	// Channel to handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Main loop to consume messages
	for {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: Terminating\n", sig)
			return

		default:
			ev := consumer.Poll(100) // Poll for events with a timeout
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Received message from topic %s [%d] at offset %v: %s\n",
					*e.TopicPartition.Topic,
					e.TopicPartition.Partition,
					e.TopicPartition.Offset,
					string(e.Value))

				processMessage(e)

			case kafka.Error:
				// Errors are informational and may not be fatal
				fmt.Printf("Consumer error: %v (%v)\n", e.Code(), e)

			default:
				fmt.Printf("Ignored event: %v\n", e)
			}
		}
	}
}
