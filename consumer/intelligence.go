package main

import (
	"fmt"
	"log"
	"logs-mgmt/db"
	"math"
	"time"
)

// DetermineLogLevel intelligently determines the log level based on various factors
func DetermineLogLevel(event *db.FlightEvent) db.LogLevel {
	logLevel := db.LogLevelInfo

	// Event Status
	switch event.EventData.EventStatus {
	case db.FlightEventStatusAccident, db.FlightEventStatusEmergency:
		return db.LogLevelFatal

	case db.FlightEventStatusCancelled:
		if logLevel < db.LogLevelError {
			logLevel = db.LogLevelError
		}

	case db.FlightEventStatusDelayed:
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}

	case db.FlightEventStatusLanded, db.FlightEventStatusArrived:
		logLevel = db.LogLevelInfo
	}

	// Battery
	if event.SystemInfo.Battery < 10.0 {
		return db.LogLevelFatal
	} else if event.SystemInfo.Battery < 20.0 {
		if logLevel < db.LogLevelError {
			logLevel = db.LogLevelError
		}
	} else if event.SystemInfo.Battery < 30.0 {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	}

	// CPU Usage
	if event.SystemInfo.CPU > 95.0 {
		if logLevel < db.LogLevelError {
			logLevel = db.LogLevelError
		}
	} else if event.SystemInfo.CPU > 85.0 {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	}

	// Memory usage
	if event.SystemInfo.Memory > 95.0 {
		if logLevel < db.LogLevelError {
			logLevel = db.LogLevelError
		}
	} else if event.SystemInfo.Memory > 85.0 {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	}

	// Disk usage
	if event.SystemInfo.Disk > 95.0 {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	}

	// Altitude (too low or too high can be dangerous)
	if event.EventData.Altitude < 50.0 && event.EventData.EventStatus != db.FlightEventStatusLanded {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	} else if event.EventData.Altitude > 15000.0 {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	}

	// Speed (unrealistic speeds)
	if event.EventData.Speed > 500.0 {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	} else if event.EventData.Speed < 50.0 && (event.EventData.EventStatus == db.FlightEventStatusDeparted || event.EventData.EventStatus == db.FlightEventStatusOnTime) {
		if logLevel < db.LogLevelWarning {
			logLevel = db.LogLevelWarning
		}
	}

	return logLevel
}

type FlightHealthStatus int32

const (
	FlightHealthGood FlightHealthStatus = iota + 1
	FlightHealthCaution
	FlightHealthWarning
	FlightHealthCritical
	FlightHealthUnknown
)

var flightHealthMap = map[FlightHealthStatus]string{
	FlightHealthGood:     "GOOD",
	FlightHealthCaution:  "CAUTION",
	FlightHealthWarning:  "WARNING",
	FlightHealthCritical: "CRITICAL",
	FlightHealthUnknown:  "UNKNOWN",
}

type FlightHealthReport struct {
	FlightID   int32              `json:"flight_id"`
	AircraftID int32              `json:"aircraft_id"`
	Health     FlightHealthStatus `json:"health"`
	Issues     []string           `json:"issues"`
	Timestamp  time.Time          `json:"time"`
}

// calculateDistance calculates the distance between two lat/lon points in kilometers
// Uses flat-earth approximation suitable for short distances (5-50km)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// 1 degree latitude ≈ 111 km
	const kmPerDegree = 111.0

	// For short distances, we can use simple Pythagorean distance
	dLat := (lat2 - lat1) * kmPerDegree
	dLon := (lon2 - lon1) * kmPerDegree

	// return sqrt(dLat*dLat + dLon*dLon)
	return math.Sqrt(dLat*dLat + dLon*dLon)
}

// JudgeFlightHealth analyzes the flight event and global data to determine flight health
func JudgeFlightHealth(event *db.FlightEvent) FlightHealthReport {
	report := FlightHealthReport{
		FlightID:   event.FlightID,
		AircraftID: event.AircraftID,
		Health:     FlightHealthGood,
		Issues:     make([]string, 0),
		Timestamp:  event.Timestamp,
	}

	// Get full flight data from sync cache
	flightData, isRunning := runningFlights.GetFlightData(event.FlightID)
	if !isRunning {
		report.Health = FlightHealthUnknown
		report.Issues = append(report.Issues, "Flight not found in active flights database")
		return report
	}

	// Check if past scheduled end time
	if flightData.EndTime.Valid && event.Timestamp.After(flightData.EndTime.Time) {
		if report.Health < FlightHealthWarning {
			report.Health = FlightHealthWarning
		}

		delay := event.Timestamp.Sub(flightData.EndTime.Time)
		report.Issues = append(report.Issues, fmt.Sprintf("Flight overdue by %.0f minutes", delay.Minutes()))
	}

	// Calculate distance to destination
	if flightData.DestinationLat != 0 && flightData.DestinationLon != 0 {
		distanceToDestination := calculateDistance(
			float64(event.EventData.Latitude),
			float64(event.EventData.Longitude),
			float64(flightData.DestinationLat),
			float64(flightData.DestinationLon),
		)

		// If very close to destination but still high altitude/speed, that's concerning
		if distanceToDestination < 5.0 && event.EventData.Altitude > 500.0 && report.Health < FlightHealthCaution {
			report.Health = FlightHealthCaution
			report.Issues = append(report.Issues, fmt.Sprintf("Close to destination (%.1fkm) but high altitude: %.0fm", distanceToDestination, event.EventData.Altitude))
		}
	}

	// Check flight duration
	if flightData.StartTime.Valid {
		duration := event.Timestamp.Sub(flightData.StartTime.Time)
		if duration.Hours() > 12.0 { // Unusually long flight
			if report.Health < FlightHealthWarning {
				report.Health = FlightHealthWarning
			}
			report.Issues = append(report.Issues, fmt.Sprintf("Unusually long flight duration: %.1f hours", duration.Hours()))
		}
	}

	// Critical issues (FlightHealthCritical)
	if event.SystemInfo.Battery < 15.0 {
		report.Health = FlightHealthCritical
		report.Issues = append(report.Issues, fmt.Sprintf("Critical battery level: %.2f%%", event.SystemInfo.Battery))
	}

	if event.EventData.EventStatus == db.FlightEventStatusAccident || event.EventData.EventStatus == db.FlightEventStatusEmergency {
		report.Health = FlightHealthCritical
		report.Issues = append(report.Issues, fmt.Sprintf("Critical flight status: %s", event.EventData.EventStatus.String()))
	}

	// Warning issues (FlightHealthWarning)
	if event.SystemInfo.CPU > 90.0 && report.Health < FlightHealthWarning {
		report.Health = FlightHealthWarning
		report.Issues = append(report.Issues, fmt.Sprintf("High CPU usage: %.2f%%", event.SystemInfo.CPU))
	}

	if event.SystemInfo.Memory > 90.0 && report.Health < FlightHealthWarning {
		report.Health = FlightHealthWarning
		report.Issues = append(report.Issues, fmt.Sprintf("High memory usage: %.2f%%", event.SystemInfo.Memory))
	}

	if event.SystemInfo.Battery < 25.0 && report.Health < FlightHealthWarning {
		report.Health = FlightHealthWarning
		report.Issues = append(report.Issues, fmt.Sprintf("Low battery: %.2f%%", event.SystemInfo.Battery))
	}

	if event.EventData.EventStatus == db.FlightEventStatusDelayed && report.Health < FlightHealthWarning {
		report.Health = FlightHealthWarning
		report.Issues = append(report.Issues, "Flight is delayed")
	}

	// Caution issues (FlightHealthCaution)
	if event.SystemInfo.CPU > 75.0 && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("Elevated CPU usage: %.2f%%", event.SystemInfo.CPU))
	}

	if event.SystemInfo.Memory > 75.0 && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("Elevated memory usage: %.2f%%", event.SystemInfo.Memory))
	}

	if event.SystemInfo.Battery < 40.0 && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("Battery level: %.2f%%", event.SystemInfo.Battery))
	}

	if event.SystemInfo.Disk > 90.0 && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("High disk usage: %.2f%%", event.SystemInfo.Disk))
	}

	// Altitude checks
	if event.EventData.Altitude < 100.0 && event.EventData.EventStatus != db.FlightEventStatusLanded && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("Low altitude: %.2fm", event.EventData.Altitude))
	}

	// Speed checks
	if event.EventData.Speed < 80.0 && (event.EventData.EventStatus == db.FlightEventStatusDeparted || event.EventData.EventStatus == db.FlightEventStatusOnTime) && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("Low speed: %.2f km/h", event.EventData.Speed))
	}

	// Network connectivity
	if event.SystemInfo.Network < 20.0 && report.Health < FlightHealthCaution {
		report.Health = FlightHealthCaution
		report.Issues = append(report.Issues, fmt.Sprintf("Weak network: %.2f Mbps", event.SystemInfo.Network))
	}

	return report
}

func LogFlightHealth(report FlightHealthReport) {
	log.Printf("%s - flight_id:%d, aircraft_id:%d", flightHealthMap[report.Health], report.FlightID, report.AircraftID)
}
