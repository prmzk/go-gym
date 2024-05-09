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
sets.updated_at as set_updated_at
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
INSERT INTO workout_exercises (id, workout_id, exercise_id, created_at) VALUES (  
  unnest(@id_array::UUID[]),  
  unnest(@workout_id_array::UUID[]),
  unnest(@exercise_id_array::UUID[]), 
  unnest(@created_at_array::TIMESTAMPTZ[])  
)
RETURNING *;

-- name: CreateSets :many
INSERT INTO sets (id, workout_exercise_id, weight, deducted_weight, duration, reps, created_at) VALUES (  
  unnest(sqlc.narg('id_array')::UUID[]),  
  unnest(sqlc.narg('workout_exercise_id_array')::UUID[]),  
  NULLIF(unnest(sqlc.narg('weight_array')::FLOAT4[]), 0),  
  NULLIF(unnest(sqlc.narg('deducted_weight_array')::FLOAT4[]), 0),  
  NULLIF(unnest(sqlc.narg('duration_array')::INT[]), 0),  
  NULLIF(unnest(sqlc.narg('reps_array')::INT[]), 0),  
  unnest(sqlc.narg('created_at_array')::TIMESTAMPTZ[])  
)
RETURNING *;

