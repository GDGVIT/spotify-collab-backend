-- name: CreateUser :one
INSERT INTO users(email, spotify_id, name)
VALUES ($1, $2, $3)
RETURNING user_uuid, id, created_at, version;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByUUID :one
SELECT *
FROM users
WHERE user_uuid = $1;

-- name: GetUserBySpotifyID :one
SELECT user_uuid
FROM users 
WHERE spotify_id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_uuid = $1;
