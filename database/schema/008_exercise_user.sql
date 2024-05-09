-- +goose Up
CREATE TABLE exercise_user (
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID REFERENCES users(id),
    exercise_id UUID REFERENCES exercises(id),
    notes TEXT,
    rest_time INT,
    PRIMARY KEY (user_id, exercise_id)
);
-- +goose Down
DROP TABLE exercise_user;