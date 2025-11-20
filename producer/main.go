package main

import (
	"context"
	"flag"
	"fmt"
	"logs-mgmt/db"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

var topic = db.KAFKA_TOPIC

const DELAY = 1 * time.Second // create n (= number of flights) log entries per DELAY interval

func generateFakeEvents(producer *kafka.Producer, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		producer.Close()
	}()

	for {
		var runningFlightsData []db.RunningFlightsData
		runningFlights.Mutex.RLock()
		runningFlightsData = runningFlights.Data
		runningFlights.Mutex.RUnlock()

		for _, flight := range runningFlightsData {
			flightEvent := db.GetRandomFlightEvent(&flight)

			producer.Produce(&kafka.Message{
				Value:          flightEvent.ToByteArray(),
				TopicPartition: kafka.TopicPartition{Topic: &topic},
			}, nil)
		}

		fmt.Printf("[%s] generated %d flight events\n", time.Now().Format("2006-01-02 15:04:05"), len(runningFlightsData))

		time.Sleep(DELAY)
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	// accept a flag to seed the database
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	if *seed {
		ctx := context.Background()
		err := db.CreateSeedData(ctx, 100)
		if err != nil {
			fmt.Printf("Failed to seed database: %s\n", err)
			os.Exit(1)
		}
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "logs-producer",
		"acks":              "all",
	})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	go syncDataPerioically(&wg)
	wg.Add(1)

	go generateFakeEvents(producer, &wg)
	wg.Add(1)

	wg.Wait()
}
