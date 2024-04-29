-- name: CreateUser :one
INSERT INTO users (id, email) 
VALUES ($1, $2) 
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: SetUserToken :one
UPDATE users SET jwt_id = $2, token_expiration = $3 WHERE id = $1
RETURNING *;

-- name: ClearUserToken :one
UPDATE users SET jwt_id = NULL, token_expiration = NULL WHERE id = $1
RETURNING *;