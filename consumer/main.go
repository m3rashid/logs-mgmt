package main

import (
	"fmt"
	"log"
	"logs-mgmt/db"
	"os"
	"os/signal"
	"sync"
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

	var wg sync.WaitGroup

	wg.Add(1)
	go syncFlightsDataPeriodically(&wg)

	eventsToMongodbBuffer := NewBufferedFlush(EventsToMongoBufferConfig)
	wg.Add(1)
	go eventsToMongodbBuffer.HandlePeriodicFlush(&wg)

	// Channel to handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Main loop to consume messages
	for {
		select {
		case sig := <-sigchan:
			log.Printf("\nCaught signal %v: Initiating graceful shutdown...\n", sig)
			eventsToMongodbBuffer.FlushBuffer()
			return

		default:
			ev := consumer.Poll(100) // Poll for events with a timeout
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				processMessage(e, eventsToMongodbBuffer)

			case kafka.Error:
				// Errors are informational and may not be fatal
				log.Printf("Consumer error: %v (%v)\n", e.Code(), e)

			default:
				// Silently ignore other event types
			}
		}
	}
}
