-- name: CreatePlaylist :one
INSERT INTO playlists (playlist_id, event_uuid, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllPlaylists :many
SELECT *
FROM playlists
WHERE event_uuid = $1;

-- name: GetPlaylist :one
SELECT *
FROM playlists
WHERE event_uuid = $1 AND playlist_id = $2;


-- name: GetPlaylistUUIDByName :one
SELECT playlist_id 
FROM playlists
WHERE event_uuid = $1 AND name = $2;

-- name: UpdatePlaylistName :one
UPDATE playlists
SET name = $1
WHERE event_uuid = $2 AND playlist_id = $3
RETURNING *;

-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE event_uuid = $1 AND playlist_id = $2;