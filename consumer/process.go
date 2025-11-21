package main

import (
	"context"
	"fmt"
	"logs-mgmt/db"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func processMessage(message *kafka.Message) {
	flightEvent, err := db.ToFlightEvent(message.Value)
	if err != nil {
		fmt.Printf("Error parsing flight event: %v\n", err)
		return
	}

	flightEvent.ID = primitive.NewObjectIDFromTimestamp(flightEvent.Timestamp)

	mongoDb, err := db.GetMongoDB()
	if err != nil {
		fmt.Printf("Error getting mongo db: %v\n", err)
		return
	}

	_, err = mongoDb.Collection(db.FLIGHT_EVENTS_COLLECTION).InsertOne(context.Background(), flightEvent)
	if err != nil {
		fmt.Printf("Error inserting flight event: %v\n", err)
		return
	}
}
