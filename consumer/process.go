package main

import (
	"context"
	"log"
	"logs-mgmt/db"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BATCH_SIZE    = 100             // Number of events to batch before insert
	BATCH_TIMEOUT = 5 * time.Second // Maximum time to wait before flushing batch
)

type EventBuffer struct {
	Events    []*db.FlightEvent
	Mutex     sync.Mutex
	LastFlush time.Time
}

var eventBuffer = EventBuffer{
	Events:    make([]*db.FlightEvent, 0, BATCH_SIZE),
	Mutex:     sync.Mutex{},
	LastFlush: time.Now(),
}

func (eb *EventBuffer) AddToBuffer(event *db.FlightEvent) {
	eb.Mutex.Lock()
	defer eb.Mutex.Unlock()

	eb.Events = append(eb.Events, event)

	if len(eb.Events) >= BATCH_SIZE || time.Since(eb.LastFlush) >= BATCH_TIMEOUT {
		eb.flushUnlocked()
	}
}

// flushUnlocked flushes the buffer (must be called with lock held)
func (eb *EventBuffer) flushUnlocked() {
	if len(eb.Events) == 0 {
		return
	}

	mongoDb, err := db.GetMongoDB()
	if err != nil {
		log.Println("Error getting mongo db", err)
		return
	}

	docs := make([]any, len(eb.Events))
	for i, event := range eb.Events {
		docs[i] = event
	}

	if _, err = mongoDb.Collection(db.FLIGHT_EVENTS_COLLECTION).InsertMany(context.Background(), docs); err != nil {
		log.Printf("Error batch inserting %d flight events: %v\n", len(eb.Events), err)
		return
	}

	eb.Events = make([]*db.FlightEvent, 0, BATCH_SIZE)
	eb.LastFlush = time.Now()

	log.Printf("Inserted %d events\n", len(eb.Events))
}

func (eb *EventBuffer) FlushBuffer() {
	eb.Mutex.Lock()
	defer eb.Mutex.Unlock()
	eb.flushUnlocked()
}

// periodicFlush flushes the buffer periodically to ensure events don't sit too long
func periodicFlush(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(BATCH_TIMEOUT)
	defer ticker.Stop()

	for range ticker.C {
		eventBuffer.FlushBuffer()
	}
}

func processMessage(message *kafka.Message) {
	flightEvent, err := db.ToFlightEvent(message.Value)
	if err != nil {
		log.Printf("Error parsing flight event: %v\n", err)
		return
	}

	flightEvent.LogLevel = DetermineLogLevel(flightEvent)
	flightEvent.ID = primitive.NewObjectIDFromTimestamp(flightEvent.Timestamp)

	healthReport := JudgeFlightHealth(flightEvent)

	if healthReport.Health != FlightHealthGood {
		LogFlightHealth(healthReport)
	}

	eventBuffer.AddToBuffer(flightEvent)
}
