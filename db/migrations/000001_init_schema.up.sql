CREATE TABLE IF NOT EXISTS aircrafts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    serial_number VARCHAR(255) NOT NULL,
    battery_health INTEGER NOT NULL,
    aircraft_status INTEGER NOT NULL,
    last_maintenance_date TIMESTAMP NOT NULL,
    maintenance_due_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS flights (
    id SERIAL PRIMARY KEY,
    aircraft_id INTEGER NOT NULL,
    source_lat REAL NOT NULL,
    source_lon REAL NOT NULL,
    destination_lat REAL NOT NULL,
    destination_lon REAL NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    flight_status INTEGER NOT NULL,
    FOREIGN KEY (aircraft_id) REFERENCES aircrafts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS db_seeded (
  id INTEGER NOT NULL DEFAULT 1 PRIMARY KEY CHECK (id = 1),
	seeded BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_flights_aircraft_id ON flights(aircraft_id);
CREATE INDEX IF NOT EXISTS idx_flights_flight_status ON flights(flight_status);
CREATE INDEX IF NOT EXISTS idx_aircrafts_aircraft_status ON aircrafts(aircraft_status);
