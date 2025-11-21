package db

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var postgresDbPool *pgxpool.Pool

func GetPostgresDB() (*pgxpool.Pool, error) {
	if postgresDbPool == nil {
		_db, err := getPostgresDbConnection()
		if err != nil {
			return nil, err
		}
		postgresDbPool = _db
	}

	return postgresDbPool, nil
}

var NO_ROWS_FOUND_ERROR_TEXT = "no rows in result set"

func getPostgresDbConnection() (*pgxpool.Pool, error) {
	postgresUrl := os.Getenv("POSTGRES_URL")
	if postgresUrl == "" {
		return nil, errors.New("POSTGRES_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func CreateSeedData(ctx context.Context, count int) error {
	pgConn, err := GetPostgresDB()
	if err != nil {
		return err
	}

	queries := New(pgConn)

	// check if the database is seeded
	seeded, err := queries.IsDbSeeded(ctx)
	if err != nil {
		if err.Error() == NO_ROWS_FOUND_ERROR_TEXT {
			fmt.Println("Database is not seeded, seeding database")
		} else {
			return err
		}
	}

	if seeded {
		fmt.Println("Database is already seeded, skipping seeding")
		return nil
	}

	fmt.Println("Seeding database...")

	// create count number of aircrafts with random data
	for range count {
		aircraftId, err := queries.CreateAircraft(ctx, CreateAircraftParams{
			Name:                gofakeit.Name(),
			Model:               gofakeit.CarModel(),
			SerialNumber:        gofakeit.Person().CreditCard.Number,
			BatteryHealth:       rand.Int31n(101),
			Status:              GetRandomStatus(),
			LastMaintenanceDate: pgtype.Timestamp{Time: time.Now(), Valid: true},
			MaintenanceDueDate:  pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 30), Valid: true}, // 30 days from now
		})
		if err != nil {
			return err
		}

		randomNumber := rand.Intn(100)
		if randomNumber%2 == 0 { // 50% chance to create a flight for the aircraft
			_, err := queries.CreateFlight(ctx, CreateFlightParams{
				AircraftID:         aircraftId,
				SourceLat:          gofakeit.Float32(),
				SourceLon:          gofakeit.Float32(),
				DestinationLat:     gofakeit.Float32(),
				DestinationLon:     gofakeit.Float32(),
				StartTime:          pgtype.Timestamp{Time: time.Now(), Valid: true},
				EndTime:            pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 30), Valid: true}, // 30 days from now
				CurrentLocationLat: gofakeit.Float32(),
				CurrentLocationLon: gofakeit.Float32(),
				Status:             GetRandomFlightStatus(),
			})
			if err != nil {
				return err
			}
		}
	}

	err = queries.SetDbSeeded(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Database seeded successfully")

	return nil
}
