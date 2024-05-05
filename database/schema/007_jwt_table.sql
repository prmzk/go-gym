-- +goose Up
ALTER TABLE users
DROP COLUMN jwt_id,
DROP COLUMN token_expiration;

CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID REFERENCES users(id),
    expiration TIMESTAMPTZ NOT NULL,
    type TEXT NOT NULL
);

-- +goose Down
ALTER TABLE users
ADD COLUMN jwt_id TEXT DEFAULT NULL UNIQUE,
ADD COLUMN token_expiration TIMESTAMPTZ DEFAULT NULL;

DROP TABLE tokens;