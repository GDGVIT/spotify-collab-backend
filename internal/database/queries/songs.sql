-- name: AddSong :one
INSERT INTO songs (song_uri, playlist_uuid)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllSongs :many
SELECT * 
FROM songs
WHERE playlist_uuid = $1;

-- name: IncreaseSongCount :one
UPDATE songs
SET count = count + 1
WHERE song_uri = $1
RETURNING count;

-- name: DecreaseSongCount :one
UPDATE songs
SET count = count - 1
WHERE song_uri = $1
RETURNING count;

-- name: DeleteSong :execrows
DELETE FROM songs
WHERE song_uri = $1;