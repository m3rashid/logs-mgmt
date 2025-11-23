package main

import (
	"context"
	"log"
	"logs-mgmt/db"
	"sync"
	"time"
)

const SYNC_INTERVAL = 2 * time.Second

type RunningFlights struct {
	Data  map[int32]*db.Flight // key: flightID, value: full flight data
	Mutex sync.RWMutex
}

var runningFlights = RunningFlights{
	Data:  make(map[int32]*db.Flight),
	Mutex: sync.RWMutex{},
}

// GetFlightData returns the complete flight data for a given flight ID
func (rf *RunningFlights) GetFlightData(flightID int32) (*db.Flight, bool) {
	rf.Mutex.RLock()
	defer rf.Mutex.RUnlock()
	data, exists := rf.Data[flightID]
	return data, exists
}

func syncFlightsDataPeriodically(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		pgConn, err := db.GetPostgresDB()
		if err != nil {
			log.Printf("Error getting postgres db: %v\n", err)
			time.Sleep(SYNC_INTERVAL)
			continue
		}

		queries := db.New(pgConn)
		flights, err := queries.GetFlights(context.Background(), db.GetFlightsParams{
			Limit:        10000, // High limit to get all in-flight flights
			Offset:       0,
			FlightStatus: db.FlightStatusInFlight,
		})
		if err != nil {
			log.Printf("Error getting in-flight flights: %v\n", err)
			time.Sleep(SYNC_INTERVAL)
			continue
		}

		runningFlights.Mutex.Lock()
		runningFlights.Data = make(map[int32]*db.Flight)
		for i := range flights {
			runningFlights.Data[flights[i].ID] = &flights[i]
		}
		runningFlights.Mutex.Unlock()

		time.Sleep(SYNC_INTERVAL)
	}
}
