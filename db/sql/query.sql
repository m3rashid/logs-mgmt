-- name: IsDbSeeded :one
SELECT seeded FROM db_seeded LIMIT 1;

-- name: SetDbSeeded :exec
INSERT INTO db_seeded (id, seeded) VALUES (1, TRUE) ON CONFLICT (id) DO UPDATE SET seeded = TRUE;

-- name: CreateAircraft :one
INSERT INTO aircrafts (
	name,
	model,
	serial_number,
	battery_health,
	status,
	last_maintenance_date,
	maintenance_due_date
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: GetAircrafts :many
SELECT * FROM aircrafts ORDER BY id DESC LIMIT $1 OFFSET $2;

-- name: GetAircraft :one
SELECT * FROM aircrafts WHERE id = $1 LIMIT 1;

-- name: GetFlights :many
SELECT * FROM flights ORDER BY id DESC LIMIT $1 OFFSET $2;

-- name: CreateFlight :one
INSERT INTO flights (
	aircraft_id,
	source_lat,
	source_lon,
	destination_lat,
	destination_lon,
	start_time,
	end_time,
	current_location_lat,
	current_location_lon,
	status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;

-- name: GetRunningFlights :many
SELECT id, aircraft_id FROM flights WHERE status = $1;
