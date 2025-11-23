-- name: GetFlight :one
SELECT * FROM flights WHERE id = $1 LIMIT 1;

-- name: CreateFlight :one
INSERT INTO flights (
	aircraft_id,
	source_lat,
	source_lon,
	destination_lat,
	destination_lon,
	start_time,
	end_time,
	flight_status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;

-- name: GetFlights :many
SELECT * FROM flights
WHERE
	flight_status = COALESCE(sqlc.narg('flight_status'), flight_status)
  AND (sqlc.narg('start_date')::timestamptz IS NULL OR start_time >= sqlc.narg('start_date'))
  AND (sqlc.narg('end_date')::timestamptz IS NULL OR end_time <= sqlc.narg('end_date'))
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateFlight :one
UPDATE flights SET
	flight_status = COALESCE($2, flight_status),
	end_time = COALESCE($3, end_time)
WHERE id = $1 RETURNING *;

-- name: GetRunningFlights :many
SELECT id, aircraft_id FROM flights WHERE flight_status = $1;
