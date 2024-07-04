-- name: CreateEvent :one
INSERT INTO events (user_uuid, event_uuid, name, event_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllEvents :many
SELECT *
FROM events
WHERE user_uuid = $1;

-- name: GetEvent :one
SELECT *
FROM events
WHERE user_uuid = $1 AND event_uuid = $2;

-- name: GetEventUUIDByName :one
SELECT event_uuid
FROM events
WHERE name = $1;

-- name: UpdateEventName :one
UPDATE events
SET name = $1
WHERE user_uuid = $2 AND event_uuid = $3
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE user_uuid = $1 AND event_uuid = $2;