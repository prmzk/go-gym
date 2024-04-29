-- +goose Up
ALTER TABLE users
ADD COLUMN jwt_id TEXT DEFAULT NULL UNIQUE,
ADD COLUMN token_expiration TIMESTAMPTZ DEFAULT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN jwt_id,
DROP COLUMN token_expiration;