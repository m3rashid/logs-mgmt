package db

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AircraftStatus int

const (
	AircraftStatusInactive AircraftStatus = iota + 1
	AircraftStatusActive
	AircraftStatusMaintenance
)

type Position = [2]float64 // [latitude, longitude]

type FlightStatus int

const (
	FlightStatusScheduled FlightStatus = iota + 1
	FlightStatusInFlight
	FlightStatusCompleted
	FlightStatusFailed
)

type LogLevel int

const (
	LogLevelNone LogLevel = iota + 1
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
)

type SystemInfo struct {
	CPU     float64 `json:"cpu" bson:"cpu"`
	Memory  float64 `json:"memory" bson:"mem"`
	Disk    float64 `json:"disk" bson:"disk"`
	Network float64 `json:"network" bson:"net"`
	Battery float64 `json:"battery" bson:"bat"`
	Uptime  float64 `json:"uptime" bson:"up"`
}

type FlightEventStatus int

const (
	FlightEventStatusEarly FlightEventStatus = iota + 1
	FlightEventStatusOnTime
	FlightEventStatusDeparted
	FlightEventStatusArrived
	FlightEventStatusLanded
	FlightEventStatusDelayed
	FlightEventStatusEmergency
	FlightEventStatusCancelled
	FlightEventStatusAccident
)

func (f FlightEventStatus) String() string {
	return []string{
		"Early",
		"On Time",
		"Departed",
		"Arrived",
		"Landed",
		"Delayed",
		"Cancelled",
		"Emergency",
		"Accident",
	}[f-1] // Subtract 1 to convert to 0-based index
}

type FlightEventData struct {
	Speed       float64           `json:"speed" bson:"speed"`
	Latitude    float64           `json:"lat" bson:"lat"`
	Longitude   float64           `json:"lon" bson:"lon"`
	Altitude    float64           `json:"alt" bson:"alt"`
	EventStatus FlightEventStatus `json:"status" bson:"status"`
}

type FlightEvent struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	FlightID   int32              `json:"flight_id" bson:"flight_id"`
	AircraftID int32              `json:"aircraft_id" bson:"aircraft_id"`
	Timestamp  time.Time          `json:"time" bson:"time"`
	LogLevel   LogLevel           `json:"level" bson:"level"`
	EventData  FlightEventData    `json:"data" bson:"evt"`
	SystemInfo SystemInfo         `json:"sys_info" bson:"sys"`
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
