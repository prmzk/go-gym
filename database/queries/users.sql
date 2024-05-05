-- name: CreateUser :one
INSERT INTO users (id, email) 
VALUES ($1, $2) 
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: SetUserToken :one
INSERT INTO tokens (user_id, expiration, type)
VALUES (sqlc.arg('user_id'), sqlc.arg('expiration'), sqlc.arg('type'))
RETURNING *;

-- name: GetUserToken :one
SELECT tokens.*, users.email as user_email, users.id as user_id, users.created_at as user_created_at, users.updated_at as user_updated_at, users.name as user_name
FROM tokens
INNER JOIN users ON tokens.user_id = users.id
WHERE tokens.id = sqlc.arg('token_id');

-- name: ClearUserToken :one
DELETE FROM tokens WHERE (id = sqlc.arg('token_id') OR user_id = sqlc.arg('user_id'))  AND type = sqlc.arg('type')
RETURNING *;

-- name: ClearAllTokenUser :many
DELETE FROM tokens WHERE user_id = sqlc.arg('user_id')
RETURNING *;

-- name: PurgeExpiredTokens :many
DELETE FROM tokens WHERE expiration < CURRENT_TIMESTAMP
RETURNING *;