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
