package main

import (
	"context"
	"log"
	"logs-mgmt/db"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v5/pgtype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var EventsToMongoBufferConfig = BufferConfig[db.FlightEvent]{
	BatchSize:    100,
	BatchTimeout: 5 * time.Second,
	FlusherFunc: func(events []*db.FlightEvent) error {
		mongoDb, err := db.GetMongoDB()
		if err != nil {
			return err
		}

		docs := make([]any, len(events))
		for i, event := range events {
			docs[i] = event
		}

		if _, err = mongoDb.Collection(db.FLIGHT_EVENTS_COLLECTION).InsertMany(context.Background(), docs); err != nil {
			log.Printf("Error batch inserting %d flight events: %v\n", len(events), err)
			return err
		}

		return nil
	},
}

type ChangeFlightStatusData struct {
	FlightId     int32
	FlightStatus db.FlightStatus
	EndTime      time.Time
}

func changeFlightStatus(flightId int32, flightStatus db.FlightStatus, endTime time.Time) {
	dbConn, err := db.GetPostgresDB()
	if err != nil {
		return
	}

	queries := db.New(dbConn)
	if _, err = queries.UpdateFlight(context.Background(), db.UpdateFlightParams{
		ID:           flightId,
		FlightStatus: flightStatus,
		EndTime:      pgtype.Timestamp{Time: endTime, Valid: true},
	}); err != nil {
		return
	}
}

func processMessage(message *kafka.Message, eventsToMongodbBuffer *BufferedFlush[db.FlightEvent]) {
	flightEvent, err := db.ToFlightEvent(message.Value)
	if err != nil {
		log.Printf("Error parsing flight event: %v\n", err)
		return
	}

	flightEvent.LogLevel = DetermineLogLevel(flightEvent)
	healthReport := JudgeFlightHealth(flightEvent)

	if healthReport.Health != FlightHealthGood {
		LogFlightHealth(healthReport)
	}

	flightEvent.ID = primitive.NewObjectIDFromTimestamp(flightEvent.Timestamp)
	eventsToMongodbBuffer.AddToBuffer(flightEvent)

	switch flightEvent.EventData.EventStatus {
	case db.FlightEventStatusAccident, db.FlightEventStatusCancelled:
		changeFlightStatus(flightEvent.FlightID, db.FlightStatusFailed, flightEvent.Timestamp)
	case db.FlightEventStatusArrived:
		changeFlightStatus(flightEvent.FlightID, db.FlightStatusCompleted, flightEvent.Timestamp)
	}
}
