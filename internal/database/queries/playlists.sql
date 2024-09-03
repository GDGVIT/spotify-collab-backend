-- name: CreatePlaylist :one
INSERT INTO playlists (playlist_id, user_uuid, name, playlist_code)
VALUES ($1, $2, $3, $4)
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

-- name: GetPlaylistIDByUUID :one
SELECT playlist_id
FROM playlists
WHERE playlist_uuid = $1;

-- name: GetPlaylistUUIDByCode :one
SELECT playlist_uuid
FROM playlists
WHERE playlist_code = $1;

-- name: UpdatePlaylistName :one
UPDATE playlists
SET name = $1
WHERE playlist_uuid = $2
RETURNING *;

-- name: DeletePlaylist :execrows
DELETE FROM playlists
WHERE playlist_uuid = $1;