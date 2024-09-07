-- -- name: CreateToken :exec
-- INSERT INTO tokens (hash, user_uuid, expiry, scope)
-- VALUES ($1, $2, $3, $4);

-- -- name: DeleteTokensForUser :exec
-- DELETE FROM tokens
-- WHERE scope=$1 AND user_uuid=$2;

-- name: NewToken :one
INSERT INTO tokens(refresh, access, user_uuid, expiry)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateToken :one
UPDATE tokens
SET refresh=$1, access=$2, expiry=$3
WHERE user_uuid=$4
RETURNING *;

-- name: GetOAuthToken :one
SELECT refresh, access, expiry, spotify_id
FROM tokens
INNER JOIN users on users.user_uuid = tokens.user_uuid
WHERE tokens.user_uuid=$1;