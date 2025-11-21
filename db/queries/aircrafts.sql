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
