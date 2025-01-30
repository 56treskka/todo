-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, updated_at) VALUES (
    gen_random_uuid (),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
) RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.id FROM refresh_tokens
INNER JOIN users
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1 AND refresh_tokens.expires_at > NOW() AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = $1;