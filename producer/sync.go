package main

import (
	"context"
	"fmt"
	"logs-mgmt/db"
	"sync"
	"time"
)

const SYNC_INTERVAL = 10 * time.Second

type RunningFlights struct {
	Data  []db.RunningFlightsData
	Mutex sync.RWMutex
}

var runningFlights = RunningFlights{
	Data:  []db.RunningFlightsData{},
	Mutex: sync.RWMutex{},
}

func syncDataPerioically(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		pgConn, err := db.GetPostgresDB()
		if err != nil {
			fmt.Println("Error getting postgres db:", err)
			continue
		}

		queries := db.New(pgConn)
		flights, err := queries.GetRunningFlights(context.Background(), db.FlightStatusInFlight)
		if err != nil {
			fmt.Println("Error getting running flights:", err)
			continue
		}

		cleanedFlights := []db.RunningFlightsData{}
		for _, flight := range flights {
			cleanedFlights = append(cleanedFlights, db.RunningFlightsData{
				FlightID:   flight.ID,
				AircraftID: flight.AircraftID,
			})
		}

		runningFlights.Mutex.Lock()
		runningFlights.Data = cleanedFlights
		runningFlights.Mutex.Unlock()

		fmt.Printf("[%s] Syncing data periodically. Total running flights: %d\n", time.Now().Format("15:04:05"), len(runningFlights.Data))
		time.Sleep(SYNC_INTERVAL)
	}
}
