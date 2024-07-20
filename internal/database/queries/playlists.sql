-- name: CreatePlaylist :one
INSERT INTO playlists (playlist_id, user_uuid, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListPlaylists :many
SELECT *
FROM playlists
WHERE user_uuid = $1;

-- name: GetPlaylist :one
SELECT *
FROM playlists
WHERE playlist_uuid = $1;

-- name: GetPlaylistUUIDByName :one
SELECT playlist_uuid
FROM playlists
WHERE user_uuid = $1 AND name = $2;

-- name: GetPlaylistUUIDByEventUUID :one
Select playlist_uuid
FROM events_playlists_mapping
WHERE event_uuid = $1;

-- name: UpdatePlaylistName :one
UPDATE playlists
SET name = $1
WHERE playlist_uuid = $2
RETURNING *;

-- name: DeletePlaylist :execrows
DELETE FROM playlists
WHERE playlist_uuid = $1;