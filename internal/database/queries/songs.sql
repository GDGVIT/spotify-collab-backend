-- name: AddSong :one
INSERT INTO songs (song_uri, playlist_uuid)
VALUES ($1, $2)
ON CONFLICT ON CONSTRAINT songs_pk
DO UPDATE SET count = count + 1
RETURNING *;

-- name: GetAllSongs :many
SELECT * 
FROM songs
WHERE playlist_uuid = $1 AND count != -1;

-- name: BlacklistSong :execrows
UPDATE songs
SET count = -1
WHERE song_uri = $1 AND playlist_uuid = $2
RETURNING *;

-- name: GetAllBlacklisted :many
SELECT *
FROM songs
WHERE playlist_uuid = $1 AND count = -1;

-- name: DeleteBlacklist :execrows
DELETE FROM songs
WHERE song_uri = $1 AND playlist_uuid = $2 and count = -1;