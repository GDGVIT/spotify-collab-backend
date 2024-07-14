-- name: CreatePlaylist :one
INSERT INTO playlists (playlist_id, user_uuid, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllPlaylists :many
SELECT *
FROM playlists
WHERE user_uuid = $1;

-- name: GetPlaylist :one
SELECT *
FROM playlists
WHERE user_uuid = $1 AND playlist_id = $2;

-- name: GetPlaylistUUIDByName :one
SELECT playlist_id 
FROM playlists
WHERE user_uuid = $1 AND name = $2;

-- name: GetPlaylistUUIDByEventUUID :one
Select playlist_id
FROM playlists
WHERE user_uuid = $1
ORDER BY created_at desc
Limit 1;

-- name: UpdatePlaylistName :one
UPDATE playlists
SET name = $1
WHERE user_uuid = $2 AND playlist_id = $3
RETURNING *;

-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE user_uuid = $1 AND playlist_id = $2;