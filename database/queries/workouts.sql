-- name: GetWorkouts :many
SELECT workouts.id, workouts.name, workouts.created_at, workouts.updated_at, workouts.start_time, workouts.end_time, users.id as user_id
FROM workouts
INNER JOIN users ON workouts.user_id = users.id
WHERE workouts.user_id = sqlc.arg('user_id');