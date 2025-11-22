-- name: GetFlights :many
SELECT * FROM flights
WHERE
  (sqlc.narg('status')::integer IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('start_date')::timestamptz IS NULL OR start_time >= sqlc.narg('start_date'))
  AND (sqlc.narg('end_date')::timestamptz IS NULL OR end_time <= sqlc.narg('end_date'))
ORDER BY id DESC
LIMIT $1 OFFSET $2;

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
	status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;

-- name: UpdateFlight :one
UPDATE flights SET
	status = COALESCE($2, status),
	end_time = COALESCE($3, end_time)
WHERE id = $1 RETURNING *;

-- name: GetRunningFlights :many
SELECT id, aircraft_id FROM flights WHERE status = $1;
