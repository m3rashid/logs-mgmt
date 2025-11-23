package controllers

import (
	"encoding/json"
	"logs-mgmt/db"
	"logs-mgmt/server/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetAircrafts(w http.ResponseWriter, r *http.Request) {
	paginationRes := utils.PaginationResponse[db.Aircraft]{}

	query := &GetAllAircraftQuery{}
	if err := paginationRes.ParseQueryWithParams(r, 50, query); err != nil {
		http.Error(w, "Incorrect Parameters", http.StatusBadRequest)
		return
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	searchParams := db.GetAircraftsParams{
		Limit:  paginationRes.Limit + 1,
		Offset: (paginationRes.CurrentPage - 1) * paginationRes.Limit,
	}

	if query.Model != "" {
		searchParams.Model = &query.Model
	}

	if query.SerialNumber != "" {
		searchParams.SerialNumber = &query.SerialNumber
	}

	if query.MinBattery != 0 {
		searchParams.MinBattery = &query.MinBattery
	}

	queries := db.New(pgConn)
	aircrafts, err := queries.GetAircrafts(r.Context(), searchParams)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	paginationRes.Docs = aircrafts
	paginationRes.BuildPaginationResponse()

	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJson(paginationRes))
}

func GetSingleAircraft(w http.ResponseWriter, r *http.Request) {
	aircraftId := r.PathValue("aircraftId")
	if aircraftId == "" {
		http.Error(w, "Aircraft ID is required", http.StatusBadRequest)
		return
	}

	id, err := utils.ParseInt32(aircraftId)
	if err != nil {
		http.Error(w, "Invalid aircraft ID", http.StatusBadRequest)
		return
	}

	pgConn, err := db.GetPostgresDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	queries := db.New(pgConn)
	aircraft, err := queries.GetAircraft(r.Context(), id)
	if err != nil {
		http.Error(w, "Aircraft not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJson(aircraft))
}

func CreateAircraft(w http.ResponseWriter, r *http.Request) {
	body := createAircraftBody{}
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

	aircraftId, err := queries.CreateAircraft(r.Context(), db.CreateAircraftParams{
		Name:                body.Name,
		Model:               body.Model,
		SerialNumber:        body.SerialNumber,
		BatteryHealth:       body.BatteryHealth,
		AircraftStatus:      body.AircraftStatus,
		LastMaintenanceDate: pgtype.Timestamp{Time: body.LastMaintenanceDate, Valid: true},
		MaintenanceDueDate:  pgtype.Timestamp{Time: body.MaintenanceDueDate, Valid: true},
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(map[string]interface{}{"id": aircraftId})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
