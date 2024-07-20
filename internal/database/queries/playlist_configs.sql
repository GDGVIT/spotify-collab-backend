-- name: CreateDefaultConfiguration :one
INSERT INTO playlist_configs (playlist_uuid)
VALUES ($1)
RETURNING *;

-- name: UpdateConfiguration :one
UPDATE playlist_configs
SET explicit = $1, require_approval = $2, max_song = $3
WHERE playlist_uuid = $4
RETURNING *;

-- name: GetConfiguration :one
SELECT * 
FROM playlist_configs
WHERE playlist_uuid = $1;
