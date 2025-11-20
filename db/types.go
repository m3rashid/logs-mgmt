package db

import (
	"encoding/json"
	"time"
)

type AircraftStatus int

const (
	AircraftStatusInactive AircraftStatus = iota
	AircraftStatusActive
	AircraftStatusMaintenance
)

type Position = [2]float64 // [latitude, longitude]

type FlightStatus int

const (
	FlightStatusScheduled FlightStatus = iota
	FlightStatusInFlight
	FlightStatusCompleted
	FlightStatusFailed
)

type LogLevel int

const (
	LogLevelNone LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
)

type SystemInfo struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Disk    float64 `json:"disk"`
	Network float64 `json:"network"`
	Battery float64 `json:"battery"`
	Uptime  float64 `json:"uptime"`
}

type FlightEventData struct {
	Speed     float64 `json:"speed" bson:"speed"`
	Latitude  float64 `json:"lat" bson:"lat"`
	Longitude float64 `json:"lon" bson:"lon"`
	Altitude  float64 `json:"alt" bson:"alt"`
}

type FlightEvent struct {
	ID         string          `json:"id" bson:"_id"`
	FlightID   int32           `json:"flight_id" bson:"flight_id"`
	AircraftID int32           `json:"aircraft_id" bson:"aircraft_id"`
	Timestamp  time.Time       `json:"time" bson:"time"`
	LogLevel   LogLevel        `json:"level" bson:"level"`
	EventData  FlightEventData `json:"data" bson:"data"`
	SystemInfo SystemInfo      `json:"sys_info" bson:"sys_info"`
}

func (e *FlightEvent) ToByteArray() []byte {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return jsonBytes
}

func ToFlightEvent(jsonBytes []byte) (*FlightEvent, error) {
	var e FlightEvent
	err := json.Unmarshal(jsonBytes, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
