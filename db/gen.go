package db

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

const KAFKA_TOPIC string = "logs-topic"

var flightPositions = flightPositionsSyncState{
	data: make(map[int32]*FlightEventData),
}

var aircraftSystemInfo = aircraftSystemInfoSyncState{
	data: make(map[int32]*SystemInfo),
}

type flightPositionsSyncState struct {
	sync.RWMutex
	data map[int32]*FlightEventData
}

// ResetFlight removes the position state for a specific flight
func (f *flightPositionsSyncState) ResetFlight(flightID int32) {
	f.Lock()
	defer f.Unlock()
	delete(f.data, flightID)
}

// State for consistent system info per aircraft ID
type aircraftSystemInfoSyncState struct {
	sync.RWMutex
	data map[int32]*SystemInfo
}

// ResetAircraft removes the system info state for a specific aircraft
func (a *aircraftSystemInfoSyncState) ResetAircraft(aircraftID int32) {
	a.Lock()
	defer a.Unlock()
	delete(a.data, aircraftID)
}

// GetRandomStatus returns a random aircraft status
func GetRandomStatus() AircraftStatus {
	statuses := []AircraftStatus{
		AircraftStatusInactive,
		AircraftStatusActive,
		AircraftStatusMaintenance,
	}
	return statuses[rand.Intn(len(statuses))]
}

// GetRandomFlightStatus returns a random flight status
func GetRandomFlightStatus() FlightStatus {
	statuses := []FlightStatus{
		FlightStatusScheduled,
		FlightStatusInFlight,
		FlightStatusCompleted,
		FlightStatusFailed,
	}
	return statuses[rand.Intn(len(statuses))]
}

// GetRandomFlightEventStatus returns a random flight event status
func GetRandomFlightEventStatus() FlightEventStatus {
	statuses := []FlightEventStatus{
		FlightEventStatusEarly,
		FlightEventStatusOnTime,
		FlightEventStatusArrived,
		FlightEventStatusDeparted,
		FlightEventStatusLanded,
		FlightEventStatusDelayed,
		FlightEventStatusCancelled,
		FlightEventStatusAccident,
		FlightEventStatusEmergency,
	}

	return statuses[rand.Intn(len(statuses))]
}

// GetRandomLogLevel returns a random log level with weighted distribution
func GetRandomLogLevel() LogLevel {
	// Weighted distribution: more Info logs, fewer Fatal logs
	weights := []int{5, 50, 25, 15, 5} // None, Info, Warning, Error, Fatal
	total := 0
	for _, w := range weights {
		total += w
	}

	r := rand.Intn(total)
	cumulative := 0

	levels := []LogLevel{
		LogLevelNone,
		LogLevelInfo,
		LogLevelWarning,
		LogLevelError,
		LogLevelFatal,
	}

	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return levels[i]
		}
	}

	return LogLevelInfo
}

// GetRandomSystemInfo generates system info with smooth transitions for a specific aircraft
func GetRandomSystemInfo(aircraftID int32, delay float64) SystemInfo {
	aircraftSystemInfo.Lock()
	defer aircraftSystemInfo.Unlock()

	state, exists := aircraftSystemInfo.data[aircraftID]
	if !exists {
		state = &SystemInfo{
			CPU:     20.0 + rand.Float64()*30.0, // 20-50%
			Memory:  30.0 + rand.Float64()*40.0, // 30-70%
			Disk:    40.0 + rand.Float64()*30.0, // 40-70%
			Network: 50.0 + rand.Float64()*50.0, // 50-100 Mbps
			Battery: 70.0 + rand.Float64()*30.0, // 70-100%
			Uptime:  0.0,
		}
		aircraftSystemInfo.data[aircraftID] = state
	}

	// Smooth transitions with small random changes
	smoothChange := func(current, min, max, maxDelta float64) float64 {
		delta := (rand.Float64()*2 - 1) * maxDelta // Random change between -maxDelta and +maxDelta
		newValue := current + delta

		// Clamp to min/max bounds
		if newValue < min {
			newValue = min
		}
		if newValue > max {
			newValue = max
		}
		return newValue
	}

	state.CPU = smoothChange(state.CPU, 5.0, 95.0, 5.0)
	state.Memory = smoothChange(state.Memory, 20.0, 90.0, 3.0)
	state.Disk = smoothChange(state.Disk, 30.0, 95.0, 0.1) // Disk changes slowly
	state.Network = smoothChange(state.Network, 10.0, 150.0, 10.0)

	// Battery decreases slowly over time
	state.Battery = math.Max(0.0, state.Battery-rand.Float64()*0.1)

	state.Uptime += delay // seconds

	return SystemInfo{
		CPU:     math.Round(state.CPU*100) / 100,
		Memory:  math.Round(state.Memory*100) / 100,
		Disk:    math.Round(state.Disk*100) / 100,
		Network: math.Round(state.Network*100) / 100,
		Battery: math.Round(state.Battery*100) / 100,
		Uptime:  math.Round(state.Uptime*100) / 100,
	}
}

// GetRandomFlightEventData generates flight event data with smooth position transitions for a specific flight
func GetRandomFlightEventData(flightID int32) FlightEventData {
	flightPositions.Lock()
	defer flightPositions.Unlock()

	state, exists := flightPositions.data[flightID]
	if !exists {
		state = &FlightEventData{
			Latitude:  gofakeit.Latitude(),
			Longitude: gofakeit.Longitude(),
			Altitude:  1000.0 + rand.Float64()*9000.0, // 1000-10000 meters
			Speed:     150.0 + rand.Float64()*200.0,   // 150-350 km/h
		}
		flightPositions.data[flightID] = state
	}

	// Speed changes gradually (acceleration/deceleration)
	speedChange := (rand.Float64()*2 - 1) * 5.0 // +/- 5 km/h change
	state.Speed += speedChange
	state.Speed = math.Max(100.0, math.Min(400.0, state.Speed)) // Clamp between 100-400 km/h

	// Aircraft typically climb/descend at 500-1500 feet per minute
	// Converting to meters per update (assuming updates every few seconds)
	altitudeChange := (rand.Float64()*2 - 1) * 15.0 // +/- 15 meters per update
	state.Altitude += altitudeChange
	state.Altitude = math.Max(100.0, math.Min(12000.0, state.Altitude)) // Clamp between 100-12000 meters

	// Assuming update every ~2 seconds, speed in km/h
	// Distance = speed * time / 3600 (convert to km)
	// 1 degree latitude ≈ 111 km, 1 degree longitude ≈ 111 km * cos(latitude)
	timeInterval := 2.0 / 3600.0 // 2 seconds in hours
	distanceKm := state.Speed * timeInterval

	// Random bearing (direction of travel) with some consistency
	// In reality, we'd want to maintain a more consistent heading, but adding slight variations
	bearing := rand.Float64() * 2 * math.Pi // Random direction

	// Calculate new position
	// Using simple flat-earth approximation for small distances
	latChange := (distanceKm / 111.0) * math.Cos(bearing)
	lonChange := (distanceKm / (111.0 * math.Cos(state.Latitude*math.Pi/180.0))) * math.Sin(bearing)

	state.Latitude += latChange
	state.Longitude += lonChange

	// Keep within valid bounds
	state.Latitude = math.Max(-90.0, math.Min(90.0, state.Latitude))
	state.Longitude = math.Max(-180.0, math.Min(180.0, state.Longitude))

	// Round values for output
	return FlightEventData{
		Speed:     math.Round(state.Speed*100) / 100,
		Latitude:  math.Round(state.Latitude*1000000) / 1000000,  // 6 decimal places
		Longitude: math.Round(state.Longitude*1000000) / 1000000, // 6 decimal places
		Altitude:  math.Round(state.Altitude*100) / 100,
		Status:    GetRandomFlightEventStatus(),
	}
}

type RunningFlightsData struct {
	FlightID   int32
	AircraftID int32
}

// GetRandomFlightEvent generates a complete flight event for a given flight
func GetRandomFlightEvent(data *RunningFlightsData, delay float64) FlightEvent {
	return FlightEvent{
		FlightID:   data.FlightID,
		AircraftID: data.AircraftID,
		Timestamp:  time.Now(),
		SystemInfo: GetRandomSystemInfo(data.AircraftID, delay),
		EventData:  GetRandomFlightEventData(data.FlightID),
	}
}
