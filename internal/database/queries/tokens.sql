-- name: CreateToken :exec
INSERT INTO tokens (hash, user_uuid, expiry, scope)
VALUES ($1, $2, $3, $4);

-- name: DeleteTokensForUser :exec
DELETE FROM tokens
WHERE scope=$1 AND user_uuid=$2;
