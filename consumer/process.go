package main

import (
	"fmt"
	"logs-mgmt/db"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func processMessage(message *kafka.Message) {
	flightEvent, err := db.ToFlightEvent(message.Value)
	if err != nil {
		fmt.Printf("Error parsing flight event: %v\n", err)
		return
	}

	fmt.Println("Flight event received:", flightEvent)
}
