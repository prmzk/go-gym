-- name: CreateUser :one
INSERT INTO users (id, name) 
VALUES ($1, $2) 
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;