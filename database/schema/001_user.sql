-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- +goose Down
DROP TABLE users;
