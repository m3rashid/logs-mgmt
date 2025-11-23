package controllers

import (
	"context"
	"encoding/json"
	"logs-mgmt/db"
	"logs-mgmt/server/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func getFlightDetailsRealtime() ([]byte, error) {
	pgConn, err := db.GetPostgresDB()
	if err != nil {
		return nil, err
	}

	queries := db.New(pgConn)
	flights, err := queries.GetRunningFlights(context.Background(), db.FlightStatusInFlight)
	if err != nil {
		return nil, err
	}

	flightIds := bson.A{}
	for _, flight := range flights {
		flightIds = append(flightIds, flight.ID)
	}

	getLatestFlightEventsAggregationPipeline := []bson.M{
		{"$match": bson.M{"flight_id": bson.M{"$in": flightIds}}},
		{"$sort": bson.M{"time": 1}},
		{"$group": bson.M{"_id": "$flight_id", "event": bson.M{"$last": "$$ROOT"}}},
		{"$replaceRoot": bson.M{"newRoot": "$event"}},
	}

	mongoConn, err := db.GetMongoDB()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	cursor, err := mongoConn.Collection(db.FLIGHT_EVENTS_COLLECTION).Aggregate(ctx, getLatestFlightEventsAggregationPipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []db.FlightEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	eventBytes, err := json.Marshal(events)
	if err != nil {
		return nil, err
	}

	return eventBytes, nil
}

func FlightEventsRealtime() http.HandlerFunc {
	flightDetailsStreamBroker := utils.NewStreamBroker()

	timeInterval := 2 * time.Second

	go func() {
		for {
			// if there are no connected clients, do nothing
			if flightDetailsStreamBroker.GetClientCount() == 0 {
				continue
			}

			time.Sleep(timeInterval)
			eventData, err := getFlightDetailsRealtime()
			if err != nil || eventData == nil {
				continue
			}

			flightDetailsStreamBroker.Notifier <- eventData
		}
	}()

	return flightDetailsStreamBroker.ServeRequestHttp
}
