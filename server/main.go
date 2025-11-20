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

	apiRouter.HandleFunc("GET /api/v1/aircrafts", http.HandlerFunc(controllers.GetAircrafts))
	apiRouter.HandleFunc("GET /api/v1/aircrafts/:id", http.HandlerFunc(controllers.GetSingleAircraft))
	apiRouter.HandleFunc("POST /api/v1/aircrafts", http.HandlerFunc(controllers.CreateAircraft))

	apiRouter.HandleFunc("GET /api/v1/flights", http.HandlerFunc(controllers.GetFlights))
	apiRouter.HandleFunc("GET /api/v1/flights/:id", http.HandlerFunc(controllers.GetSingleFlight))
	apiRouter.HandleFunc("POST /api/v1/flights", http.HandlerFunc(controllers.CreateFlight))
	apiRouter.HandleFunc("PUT /api/v1/flights/:id", http.HandlerFunc(controllers.UpdateFlight))

	apiRouter.HandleFunc("GET /api/v1/flights/:id/realtime", http.HandlerFunc(controllers.FlightDetailsRealtime))

	PORT := "0.0.0.0:" + os.Getenv("PORT")
	server := http.Server{
		Addr:    PORT,
		Handler: controllers.MiddlewareStack(apiRouter),
	}

	log.Printf("Server is running at http://localhost%s", PORT)
	server.ListenAndServe()
}
