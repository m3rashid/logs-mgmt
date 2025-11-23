package controllers

import (
	"encoding/json"
	"log"
	"logs-mgmt/db"
	"logs-mgmt/server/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetFlights(w http.ResponseWriter, r *http.Request) {
	paginationRes := utils.PaginationResponse[db.Flight]{}

	query := &GetAllFlightQuery{}
	if err := paginationRes.ParseQueryWithParams(r, 50, query); err != nil {
		http.Error(w, "Incorrect Parameters", http.StatusBadRequest)
		return
	}

	searchParams := db.GetFlightsParams{
		Limit:  paginationRes.Limit + 1,
		Offset: (paginationRes.CurrentPage - 1) * paginationRes.Limit,
	}

	if query.FlightStatus != 0 {
		searchParams.FlightStatus = query.FlightStatus
	}

	if query.StartDate != "" {
		startDate, err := utils.ParseTimestamp(query.StartDate)
		if err == nil {
			searchParams.StartDate = startDate
		}
	}

	if query.EndDate != "" {
		endDate, err := utils.ParseTimestamp(query.EndDate)
		if err == nil {
			searchParams.EndDate = endDate
		}
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(pgConn)
	flights, err := queries.GetFlights(r.Context(), searchParams)
	if err != nil {
		log.Println("Error getting flights", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	paginationRes.Docs = flights
	paginationRes.BuildPaginationResponse()

	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJson(paginationRes))
}

func GetSingleFlight(w http.ResponseWriter, r *http.Request) {
	flightId := r.PathValue("flightId")
	if flightId == "" {
		http.Error(w, "Flight ID is required", http.StatusBadRequest)
		return
	}

	id, err := utils.ParseInt32(flightId)
	if err != nil {
		http.Error(w, "Invalid flight ID", http.StatusBadRequest)
		return
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(pgConn)
	flight, err := queries.GetFlight(r.Context(), id)
	if err != nil {
		http.Error(w, "Flight not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJson(flight))
}

func CreateFlight(w http.ResponseWriter, r *http.Request) {
	body := createFlightBody{}
	if err := utils.ParseBodyAndValidate(w, r, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(pgConn)
	flightId, err := queries.CreateFlight(r.Context(), db.CreateFlightParams{
		AircraftID:     body.AircraftID,
		SourceLat:      body.SourceLat,
		SourceLon:      body.SourceLon,
		DestinationLat: body.DestinationLat,
		DestinationLon: body.DestinationLon,
		StartTime:      pgtype.Timestamp{Time: body.StartTime, Valid: true},
		EndTime:        pgtype.Timestamp{Time: body.EndTime, Valid: true},
		FlightStatus:   db.FlightStatusScheduled,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"id": flightId})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateFlight(w http.ResponseWriter, r *http.Request) {
	flightId := r.PathValue("flightId")
	if flightId == "" {
		http.Error(w, "Flight ID is required", http.StatusBadRequest)
		return
	}

	id, err := utils.ParseInt32(flightId)
	if err != nil {
		http.Error(w, "Invalid flight ID", http.StatusBadRequest)
		return
	}

	body := updateFlightBody{}
	if err := utils.ParseBodyAndValidate(w, r, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(pgConn)
	flight, err := queries.UpdateFlight(r.Context(), db.UpdateFlightParams{
		ID:           id,
		FlightStatus: body.FlightStatus,
		EndTime:      pgtype.Timestamp{Time: body.EndTime, Valid: true},
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJson(flight))
}
