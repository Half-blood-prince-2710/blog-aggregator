-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT id,created_at,updated_at,name FROM users WHERE name = $1;

-- name: DeleteAllUser :exec 
DELETE FROM users;

-- name: GetUsers :many
SELECT id,created_at,updated_at,name FROM users;

--name: GetUserById :one
SELECT name FROM users WHERE id = $1;