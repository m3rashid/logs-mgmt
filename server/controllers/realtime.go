package controllers

import (
	"context"
	"logs-mgmt/db"
	"logs-mgmt/server/utils"
	"net/http"
	"time"

	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func getFlightDetailsRealtime() ([]byte, error) {
	log.Println("Getting flight details realtime")

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		return nil, err
	}

	queries := db.New(pgConn)
	flights, err := queries.GetRunningFlights(context.Background(), db.FlightStatusInFlight)
	if err != nil {
		return nil, err
	}

	// get the realtime data for each flight
	flightIds := []int32{}
	for _, flight := range flights {
		flightIds = append(flightIds, flight.ID)
	}

	mongoConn, err := db.GetMongoDB()
	if err != nil {
		return nil, err
	}

	collection := mongoConn.Collection(db.FLIGHT_EVENTS_COLLECTION)
	getLatestFlightEventsAggregationPipeline := []bson.M{
		{
			"$match": bson.M{
				"flight_id": bson.M{"$in": bson.A{flightIds}},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$flight_id",
				"event": bson.M{"$last": "$$ROOT"},
			},
		},
		{
			"$project": bson.M{
				"_id":   0,
				"event": "$event",
			},
		},
	}

	ctx := context.Background()
	cursor, err := collection.Aggregate(ctx, getLatestFlightEventsAggregationPipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []db.FlightEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	eventBytes := []byte{}
	for _, event := range events {
		eventBytes = append(eventBytes, event.ToByteArray()...)
	}

	return eventBytes, nil
}

func FlightDetailsRealtime() http.HandlerFunc {
	flightDetailsStreamBroker := utils.NewStreamBroker()

	timeInterval := 1 * time.Second

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
