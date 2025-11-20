package controllers

import (
	"encoding/json"
	"logs-mgmt/db"
	"logs-mgmt/server/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetFlights(w http.ResponseWriter, r *http.Request) {
	paginationRes := utils.PaginationResponse[db.Flight]{}
	if err := paginationRes.ParseQuery(r, 50); err != nil {
		http.Error(w, "Incorrect Parameters", http.StatusBadRequest)
		return
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(pgConn)
	flights, err := queries.GetFlights(r.Context(), db.GetFlightsParams{
		Limit:  paginationRes.Limit + 1,
		Offset: (paginationRes.CurrentPage - 1) * paginationRes.Limit,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	paginationRes.Docs = flights
	paginationRes.BuildPaginationResponse()

	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJson(paginationRes))
}

func GetSingleFlight(w http.ResponseWriter, r *http.Request) {}

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

	w.Write(utils.ToJson(map[string]interface{}{"id": flightId}))
}

func UpdateFlight(w http.ResponseWriter, r *http.Request) {}

func FlightDetailsRealtime(w http.ResponseWriter, r *http.Request) {}
