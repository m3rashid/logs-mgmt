-- name: CreateAircraft :one
INSERT INTO aircrafts (
	name,
	model,
	serial_number,
	battery_health,
	aircraft_status,
	last_maintenance_date,
	maintenance_due_date
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: GetAircrafts :many
SELECT * FROM aircrafts
WHERE
  (sqlc.narg('model')::text IS NULL OR model = sqlc.narg('model'))
  AND (sqlc.narg('serial_number')::text IS NULL OR serial_number = sqlc.narg('serial_number'))
  AND (sqlc.narg('min_battery')::integer IS NULL OR battery_health >= sqlc.narg('min_battery'))
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: GetAircraft :one
SELECT * FROM aircrafts WHERE id = $1 LIMIT 1;
