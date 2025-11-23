package controllers

import (
	"logs-mgmt/db"
	"logs-mgmt/server/utils"
	"time"
)

type GetAllFlightQuery struct {
	utils.PaginationQuery
	FlightStatus db.FlightStatus `query:"status"`
	StartDate    string          `query:"start_date"`
	EndDate      string          `query:"end_date"`
}

func (q *GetAllFlightQuery) GetPagination() utils.PaginationQuery {
	return q.PaginationQuery
}

type GetAllAircraftQuery struct {
	utils.PaginationQuery
	Model        string `query:"model"`
	SerialNumber string `query:"serial_number"`
	MinBattery   int32  `query:"min_battery"`
}

func (q *GetAllAircraftQuery) GetPagination() utils.PaginationQuery {
	return q.PaginationQuery
}

type createAircraftBody struct {
	Name                string            `json:"name" validate:"required"`
	Model               string            `json:"model" validate:"required"`
	SerialNumber        string            `json:"serial_number" validate:"required"`
	BatteryHealth       int32             `json:"battery_health" validate:"required,min=0,max=100"`
	AircraftStatus      db.AircraftStatus `json:"aircraft_status" validate:"required"`
	LastMaintenanceDate time.Time         `json:"last_maintenance_date" validate:"required"`
	MaintenanceDueDate  time.Time         `json:"maintenance_due_date" validate:"required"`
}

type createFlightBody struct {
	AircraftID     int32     `json:"aircraft_id" validate:"required"`
	SourceLat      float32   `json:"source_lat" validate:"required"`
	SourceLon      float32   `json:"source_lon" validate:"required"`
	DestinationLat float32   `json:"destination_lat" validate:"required"`
	DestinationLon float32   `json:"destination_lon" validate:"required"`
	StartTime      time.Time `json:"start_time" validate:"required"`
	EndTime        time.Time `json:"end_time" validate:"required"`
}

type updateFlightBody struct {
	FlightStatus db.FlightStatus `json:"flight_status"`
	EndTime      time.Time       `json:"end_time"`
}
