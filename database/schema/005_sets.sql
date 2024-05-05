-- +goose Up
CREATE TABLE sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exercise_id UUID REFERENCES workout_exercises(id),
    weight DECIMAL,
    deducted_weight DECIMAL,
    duration INTERVAL,
    reps INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE set_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_exercise_id UUID REFERENCES template_exercises(id),
    weight DECIMAL,
    deducted_weight DECIMAL,
    duration INTERVAL,
    reps INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE sets;
DROP TABLE set_templates;
