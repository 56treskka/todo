-- name: CreateUser :one
INSERT INTO users(id, name, email, password, created_at, updated_at) VALUES (
    gen_random_uuid (),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;