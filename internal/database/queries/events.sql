-- name: CreateEvent :one
INSERT INTO events (user_uuid, name, event_code)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllEvents :many
SELECT *
FROM events
WHERE user_uuid = $1;

-- name: GetEvent :one
SELECT *
FROM events
WHERE event_uuid = $1;

-- name: GetEventUUIDByName :one
SELECT event_uuid
FROM events
WHERE name = $1;

-- name: GetEventUUIDByCode :one
Select event_uuid
FROM events
WHERE event_code = $1;

-- name: UpdateEventName :one
UPDATE events
SET name = $1
WHERE event_uuid = $2
RETURNING *;

-- name: DeleteEvent :execrows
DELETE FROM events
WHERE event_uuid = $1;