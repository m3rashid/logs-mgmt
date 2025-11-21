package main

import (
	"fmt"
	"log"
	"logs-mgmt/server/controllers"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	apiRouter := http.NewServeMux()

	apiRouter.HandleFunc("GET /api/aircrafts", controllers.GetAircrafts)
	apiRouter.HandleFunc("POST /api/aircrafts", controllers.CreateAircraft)
	apiRouter.HandleFunc("GET /api/aircrafts/{aircraftId}", controllers.GetSingleAircraft)

	apiRouter.HandleFunc("GET /api/flights", controllers.GetFlights)
	apiRouter.HandleFunc("POST /api/flights", controllers.CreateFlight)
	apiRouter.HandleFunc("GET /api/flights/realtime", controllers.FlightDetailsRealtime())
	apiRouter.HandleFunc("PUT /api/flights/{flightId}", controllers.UpdateFlight)
	apiRouter.HandleFunc("GET /api/flights/{flightId}", controllers.GetSingleFlight)

	PORT := "0.0.0.0:" + os.Getenv("PORT")
	server := http.Server{
		Addr:    PORT,
		Handler: controllers.MiddlewareStack(apiRouter),
	}

	log.Printf("Server is running at http://localhost%s", PORT)
	server.ListenAndServe()
}
