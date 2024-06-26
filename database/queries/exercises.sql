-- name: GetExercises :many
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, category.name as category, body_part.name as body_part, exercise_user.notes as notes, exercise_user.rest_time as rest_time
FROM exercises
LEFT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
LEFT JOIN exercise_user ON exercises.id = exercise_user.exercise_id AND exercise_user.user_id = sqlc.arg('user_id')
WHERE (exercises.name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
AND (category.name = sqlc.narg('category') OR sqlc.narg('category') IS NULL)
AND (body_part.name = sqlc.narg('body_part') OR sqlc.narg('body_part') IS NULL)
ORDER BY exercises.name ASC;

-- name: GetExerciseById :one
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, category.name as category, body_part.name as body_part, exercise_user.notes as notes, exercise_user.rest_time as rest_time
FROM exercises
LEFT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
LEFT JOIN exercise_user ON exercises.id = exercise_user.exercise_id AND exercise_user.user_id = sqlc.arg('user_id')
WHERE exercises.id = $1;

-- name: GetExerciseByCategory :many
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, body_part.name as body_part, category.name as category_name, category.id as category_id, category.created_at as category_created_at, category.updated_at as category_updated_at, exercise_user.notes as notes, exercise_user.rest_time as rest_time
FROM exercises
RIGHT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
LEFT JOIN exercise_user ON exercises.id = exercise_user.exercise_id AND exercise_user.user_id = sqlc.arg('user_id')
WHERE category.id = sqlc.arg('category_id');

-- name: GetExerciseByBodyPart :many
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, category.name as category, body_part.name as body_part_name, body_part.id as body_part_id, body_part.created_at as body_part_created_at, body_part.updated_at as body_part_updated_at, exercise_user.notes as notes, exercise_user.rest_time as rest_time
FROM exercises
RIGHT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
LEFT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_user ON exercises.id = exercise_user.exercise_id AND exercise_user.user_id = sqlc.arg('user_id')
WHERE body_part.id = sqlc.arg('body_part_id');

-- name: GetExerciseCategories :many
SELECT * FROM exercise_categories;

-- name: GetExerciseBodyParts :many
SELECT * FROM exercise_body_parts;

-- name: UpsertExerciseUser :exec
INSERT INTO exercise_user (user_id, exercise_id, notes, rest_time)
VALUES (sqlc.arg('user_id'), sqlc.arg('exercise_id'), sqlc.narg('notes'), sqlc.narg('rest_time'))
ON CONFLICT (user_id, exercise_id) 
DO UPDATE SET notes = COALESCE(sqlc.arg('notes'), exercise_user.notes), rest_time = COALESCE(sqlc.arg('rest_time'), exercise_user.rest_time);