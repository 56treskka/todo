-- name: CreateTask :one
INSERT INTO tasks (id, title, description, user_id, created_at, updated_at) VALUES (gen_random_uuid (), $1, $2, $3, NOW(), NOW()) RETURNING *;

-- name: UpdateTask :one
UPDATE tasks SET title = $1, description = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4 RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1 AND user_id = $2;

-- name: GetTasks :many
SELECT * FROM tasks WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;