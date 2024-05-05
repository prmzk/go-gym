-- +goose Up
ALTER TABLE exercises
ADD COLUMN image_url TEXT DEFAULT NULL,
ADD COLUMN video_url TEXT DEFAULT NULL;

ALTER TABLE workouts
ADD COLUMN start_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN end_time TIMESTAMPTZ DEFAULT NULL;

-- +goose Down
ALTER TABLE exercises
DROP COLUMN image_url,
DROP COLUMN video_url;

ALTER TABLE workouts
DROP COLUMN start_time,
DROP COLUMN end_time;