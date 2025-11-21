-- name: IsDbSeeded :one
SELECT seeded FROM db_seeded LIMIT 1;

-- name: SetDbSeeded :exec
INSERT INTO db_seeded (id, seeded) VALUES (1, TRUE) ON CONFLICT (id) DO UPDATE SET seeded = TRUE;
