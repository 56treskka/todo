-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, updated_at) VALUES (
    gen_random_uuid (),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
) RETURNING *;