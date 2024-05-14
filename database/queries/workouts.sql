-- name: GetWorkouts :many
SELECT workouts.id, workouts.name, workouts.created_at, workouts.updated_at, workouts.start_time, workouts.end_time, users.id as user_id
FROM workouts
INNER JOIN users ON workouts.user_id = users.id
WHERE workouts.user_id = sqlc.arg('user_id');

-- name: GetWorkoutById :many
SELECT workouts.id as id,
user_id,
workouts.name as name,
workouts.created_at as created_at,
workouts.updated_at as updated_at,
start_time,
end_time,
workout_exercises.id as workout_exercise_id,
workout_exercises.created_at as workout_exercise_created_at,
workout_exercises.updated_at as workout_exercise_updated_at,
workout_exercises.order_no as workout_exercise_order_no,
exercises.id as exercise_id,
exercises.name as exercise_name,
exercise_categories.name as category_name,
exercise_body_parts.name as body_part_name,
sets.id as set_id,
weight,
deducted_weight,
duration,
reps,
sets.created_at as set_created_at,
sets.updated_at as set_updated_at,
sets.order_no as set_order_no
FROM workouts
INNER JOIN users ON workouts.user_id = users.id
INNER JOIN workout_exercises ON workouts.id = workout_exercises.workout_id
INNER JOIN exercises ON workout_exercises.exercise_id = exercises.id
INNER JOIN exercise_categories ON exercises.category_id = exercise_categories.id
INNER JOIN exercise_body_parts ON exercises.body_part_id = exercise_body_parts.id
INNER JOIN sets ON workout_exercises.id = sets.workout_exercise_id
WHERE workouts.id = sqlc.arg('id') AND workouts.user_id = sqlc.arg('user_id');


-- name: CreateWorkout :one
INSERT INTO workouts (name, user_id, start_time, end_time, created_at) VALUES (sqlc.arg('name'), sqlc.arg('user_id'), sqlc.arg('start_time'), sqlc.arg('end_time'), sqlc.arg('created_at'))
RETURNING *;

-- name: CreateWorkoutExercise :many
INSERT INTO workout_exercises (id, workout_id, exercise_id, created_at, order_no) VALUES (  
  unnest(@id_array::UUID[]),  
  unnest(@workout_id_array::UUID[]),
  unnest(@exercise_id_array::UUID[]), 
  unnest(@created_at_array::TIMESTAMPTZ[]),
  unnest(@order_no_array::INT[])
)
RETURNING *;

-- name: CreateSets :many
INSERT INTO sets (id, workout_exercise_id, weight, deducted_weight, duration, reps, created_at, order_no) VALUES (  
  unnest(sqlc.narg('id_array')::UUID[]),  
  unnest(sqlc.narg('workout_exercise_id_array')::UUID[]),  
  NULLIF(unnest(sqlc.narg('weight_array')::FLOAT4[]), 0),  
  NULLIF(unnest(sqlc.narg('deducted_weight_array')::FLOAT4[]), 0),  
  NULLIF(unnest(sqlc.narg('duration_array')::INT[]), 0),  
  NULLIF(unnest(sqlc.narg('reps_array')::INT[]), 0),  
  unnest(sqlc.narg('created_at_array')::TIMESTAMPTZ[]),
  unnest(@order_no_array::INT[])
)
RETURNING *;

-- name: GetPreviousWorkoutExerciseSets :many
SELECT sets.id, sets.workout_exercise_id, sets.weight, sets.deducted_weight, sets.duration, sets.reps, sets.created_at, sets.order_no
FROM sets
INNER JOIN workout_exercises ON sets.workout_exercise_id = workout_exercises.id
INNER JOIN workouts ON workout_exercises.workout_id = workouts.id
WHERE workouts.user_id = sqlc.arg('user_id')
AND workout_exercises.exercise_id = sqlc.arg('exercise_id')
AND sets.workout_exercise_id IN (
  SELECT workout_exercise_id
  FROM sets
  INNER JOIN workout_exercises ON sets.workout_exercise_id = workout_exercises.id
  INNER JOIN workouts ON workout_exercises.workout_id = workouts.id
  WHERE workouts.user_id = sqlc.arg('user_id')
  AND workout_exercises.exercise_id = sqlc.arg('exercise_id')
  ORDER BY sets.created_at DESC
  LIMIT 1
)
ORDER BY sets.created_at DESC;

-- name: DeleteWorkout :exec
WITH deleted_sets AS (
  DELETE FROM sets
  WHERE workout_exercise_id IN (
    SELECT id FROM workout_exercises WHERE workout_id = sqlc.arg('workout_id')
  )
),
deleted_workout_exercises AS (
  DELETE FROM workout_exercises
  WHERE workout_id = sqlc.arg('workout_id')
)
DELETE FROM workouts
WHERE workouts.id = sqlc.arg('workout_id') AND workouts.user_id = sqlc.arg('user_id');

-- name: CreateTemplate :one
INSERT INTO templates (name, user_id) VALUES (sqlc.arg('name'), sqlc.arg('user_id'))
RETURNING *;

-- name: CreateTemplateExercises :many
INSERT INTO template_exercises (id, template_id, exercise_id, order_no) VALUES (  
  unnest(@id_array::UUID[]),  
  unnest(@template_id::UUID[]),
  unnest(@exercise_id_array::UUID[]),
  unnest(@order_no_array::INT[])
)
RETURNING *;

-- name: CreateSetTemplates :many
INSERT INTO set_templates (id, template_exercise_id, weight, deducted_weight, duration, reps, order_no) VALUES (  
  unnest(sqlc.narg('id_array')::UUID[]),  
  unnest(sqlc.narg('template_exercise_id')::UUID[]),  
  NULLIF(unnest(sqlc.narg('weight_array')::FLOAT4[]), 0),  
  NULLIF(unnest(sqlc.narg('deducted_weight_array')::FLOAT4[]), 0),  
  NULLIF(unnest(sqlc.narg('duration_array')::INT[]), 0),  
  NULLIF(unnest(sqlc.narg('reps_array')::INT[]), 0),
  unnest(@order_no_array::INT[])
)
RETURNING *;

-- name: GetTemplates :many
SELECT templates.id, templates.name, templates.created_at, templates.updated_at, users.id as user_id
FROM templates
INNER JOIN users ON templates.user_id = users.id
WHERE templates.user_id = sqlc.arg('user_id');

-- name: GetTemplateById :many
SELECT templates.id as id,
user_id,
templates.name as name,
templates.created_at as created_at,
templates.updated_at as updated_at,
template_exercises.id as template_exercise_id,
template_exercises.created_at as workout_exercise_created_at,
template_exercises.updated_at as workout_exercise_updated_at,
template_exercises.order_no as template_exercise_order_no,
exercises.id as exercise_id,
exercises.name as exercise_name,
exercise_categories.name as category_name,
exercise_body_parts.name as body_part_name,
set_templates.id as set_id,
set_templates.order_no as set_templates_order_no,
weight,
deducted_weight,
duration,
reps,
set_templates.created_at as set_created_at,
set_templates.updated_at as set_updated_at
FROM templates
INNER JOIN users ON templates.user_id = users.id
INNER JOIN template_exercises ON templates.id = template_exercises.template_id
INNER JOIN exercises ON template_exercises.exercise_id = exercises.id
INNER JOIN exercise_categories ON exercises.category_id = exercise_categories.id
INNER JOIN exercise_body_parts ON exercises.body_part_id = exercise_body_parts.id
LEFT JOIN set_templates ON template_exercises.id = set_templates.template_exercise_id
WHERE templates.id = sqlc.arg('id') AND templates.user_id = sqlc.arg('user_id');

-- name: DeleteTemplate :exec
WITH deleted_sets AS (
  DELETE FROM set_templates
  WHERE template_exercise_id IN (
    SELECT id FROM template_exercises WHERE template_id = sqlc.arg('template_id')
  )
),
deleted_template_exercises AS (
  DELETE FROM template_exercises
  WHERE template_id = sqlc.arg('template_id')
)
DELETE FROM templates
WHERE templates.id = sqlc.arg('template_id') AND templates.user_id = sqlc.arg('user_id');

-- name: DeleteTemplateExercise :exec
WITH deleted_sets AS (
  DELETE FROM set_templates
  WHERE template_exercise_id IN (
    SELECT id FROM template_exercises WHERE template_id = sqlc.arg('template_id')
  )
)
DELETE FROM template_exercises
WHERE template_exercises.template_id IN (
  SELECT id FROM templates WHERE templates.id = sqlc.arg('template_id') AND templates.user_id = sqlc.arg('user_id')
);

-- name: GetTemplateExerciseByExerciseId :many
SELECT * FROM template_exercises
WHERE template_exercises.exercise_id = ANY(sqlc.arg('exercise_id')::UUID[]);

-- name: DeleteSetTemplate :exec
DELETE FROM set_templates
WHERE set_templates.template_exercise_id IN (
  SELECT id FROM template_exercises WHERE template_exercises.template_id IN (
    SELECT id FROM templates WHERE templates.id = sqlc.arg('template_id') AND templates.user_id = sqlc.arg('user_id')
  )
);