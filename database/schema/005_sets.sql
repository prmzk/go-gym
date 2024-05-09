-- +goose Up
CREATE TABLE sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exercise_id UUID REFERENCES workout_exercises(id),
    weight FLOAT4,
    deducted_weight FLOAT4,
    duration INT,
    reps INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE set_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_exercise_id UUID REFERENCES template_exercises(id),
    weight FLOAT4,
    deducted_weight FLOAT4,
    duration INT,
    reps INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE sets;
DROP TABLE set_templates;
