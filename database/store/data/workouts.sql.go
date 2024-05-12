// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: workouts.sql

package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createSets = `-- name: CreateSets :many
INSERT INTO sets (id, workout_exercise_id, weight, deducted_weight, duration, reps, created_at) VALUES (  
  unnest($1::UUID[]),  
  unnest($2::UUID[]),  
  NULLIF(unnest($3::FLOAT4[]), 0),  
  NULLIF(unnest($4::FLOAT4[]), 0),  
  NULLIF(unnest($5::INT[]), 0),  
  NULLIF(unnest($6::INT[]), 0),  
  unnest($7::TIMESTAMPTZ[])  
)
RETURNING id, workout_exercise_id, weight, deducted_weight, duration, reps, created_at, updated_at
`

type CreateSetsParams struct {
	IDArray                []uuid.UUID
	WorkoutExerciseIDArray []uuid.UUID
	WeightArray            []float32
	DeductedWeightArray    []float32
	DurationArray          []int32
	RepsArray              []int32
	CreatedAtArray         []time.Time
}

func (q *Queries) CreateSets(ctx context.Context, arg CreateSetsParams) ([]Set, error) {
	rows, err := q.db.QueryContext(ctx, createSets,
		pq.Array(arg.IDArray),
		pq.Array(arg.WorkoutExerciseIDArray),
		pq.Array(arg.WeightArray),
		pq.Array(arg.DeductedWeightArray),
		pq.Array(arg.DurationArray),
		pq.Array(arg.RepsArray),
		pq.Array(arg.CreatedAtArray),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Set
	for rows.Next() {
		var i Set
		if err := rows.Scan(
			&i.ID,
			&i.WorkoutExerciseID,
			&i.Weight,
			&i.DeductedWeight,
			&i.Duration,
			&i.Reps,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createWorkout = `-- name: CreateWorkout :one
INSERT INTO workouts (name, user_id, start_time, end_time, created_at) VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, created_at, updated_at, start_time, end_time
`

type CreateWorkoutParams struct {
	Name      string
	UserID    uuid.NullUUID
	StartTime sql.NullTime
	EndTime   sql.NullTime
	CreatedAt time.Time
}

func (q *Queries) CreateWorkout(ctx context.Context, arg CreateWorkoutParams) (Workout, error) {
	row := q.db.QueryRowContext(ctx, createWorkout,
		arg.Name,
		arg.UserID,
		arg.StartTime,
		arg.EndTime,
		arg.CreatedAt,
	)
	var i Workout
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.StartTime,
		&i.EndTime,
	)
	return i, err
}

const createWorkoutExercise = `-- name: CreateWorkoutExercise :many
INSERT INTO workout_exercises (id, workout_id, exercise_id, created_at) VALUES (  
  unnest($1::UUID[]),  
  unnest($2::UUID[]),
  unnest($3::UUID[]), 
  unnest($4::TIMESTAMPTZ[])  
)
RETURNING id, workout_id, exercise_id, created_at, updated_at
`

type CreateWorkoutExerciseParams struct {
	IDArray         []uuid.UUID
	WorkoutIDArray  []uuid.UUID
	ExerciseIDArray []uuid.UUID
	CreatedAtArray  []time.Time
}

func (q *Queries) CreateWorkoutExercise(ctx context.Context, arg CreateWorkoutExerciseParams) ([]WorkoutExercise, error) {
	rows, err := q.db.QueryContext(ctx, createWorkoutExercise,
		pq.Array(arg.IDArray),
		pq.Array(arg.WorkoutIDArray),
		pq.Array(arg.ExerciseIDArray),
		pq.Array(arg.CreatedAtArray),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WorkoutExercise
	for rows.Next() {
		var i WorkoutExercise
		if err := rows.Scan(
			&i.ID,
			&i.WorkoutID,
			&i.ExerciseID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPreviousWorkoutExerciseSets = `-- name: GetPreviousWorkoutExerciseSets :many
WITH workout_id AS (
  SELECT workouts.id as max_date
  FROM workouts
  WHERE workouts.user_id = $2
  ORDER BY workouts.created_at DESC
  LIMIT 1
)
SELECT sets.id, sets.workout_exercise_id, sets.weight, sets.deducted_weight, sets.duration, sets.reps, sets.created_at
FROM sets
INNER JOIN workout_exercises ON sets.workout_exercise_id = workout_exercises.id
INNER JOIN workouts ON workout_exercises.workout_id = workouts.id
WHERE workouts.id = (SELECT max_date FROM workout_id)
AND workout_exercises.exercise_id = $1
`

type GetPreviousWorkoutExerciseSetsParams struct {
	ExerciseID uuid.NullUUID
	UserID     uuid.NullUUID
}

type GetPreviousWorkoutExerciseSetsRow struct {
	ID                uuid.UUID
	WorkoutExerciseID uuid.NullUUID
	Weight            sql.NullFloat64
	DeductedWeight    sql.NullFloat64
	Duration          sql.NullInt32
	Reps              sql.NullInt32
	CreatedAt         time.Time
}

func (q *Queries) GetPreviousWorkoutExerciseSets(ctx context.Context, arg GetPreviousWorkoutExerciseSetsParams) ([]GetPreviousWorkoutExerciseSetsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPreviousWorkoutExerciseSets, arg.ExerciseID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPreviousWorkoutExerciseSetsRow
	for rows.Next() {
		var i GetPreviousWorkoutExerciseSetsRow
		if err := rows.Scan(
			&i.ID,
			&i.WorkoutExerciseID,
			&i.Weight,
			&i.DeductedWeight,
			&i.Duration,
			&i.Reps,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkoutById = `-- name: GetWorkoutById :many
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
WHERE workouts.id = $1 AND workouts.user_id = $2
`

type GetWorkoutByIdParams struct {
	ID     uuid.UUID
	UserID uuid.NullUUID
}

type GetWorkoutByIdRow struct {
	ID                       uuid.UUID
	UserID                   uuid.NullUUID
	Name                     string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	StartTime                sql.NullTime
	EndTime                  sql.NullTime
	WorkoutExerciseID        uuid.UUID
	WorkoutExerciseCreatedAt time.Time
	WorkoutExerciseUpdatedAt time.Time
	ExerciseID               uuid.UUID
	ExerciseName             string
	CategoryName             string
	BodyPartName             string
	SetID                    uuid.UUID
	Weight                   sql.NullFloat64
	DeductedWeight           sql.NullFloat64
	Duration                 sql.NullInt32
	Reps                     sql.NullInt32
	SetCreatedAt             time.Time
	SetUpdatedAt             time.Time
}

func (q *Queries) GetWorkoutById(ctx context.Context, arg GetWorkoutByIdParams) ([]GetWorkoutByIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkoutById, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkoutByIdRow
	for rows.Next() {
		var i GetWorkoutByIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.StartTime,
			&i.EndTime,
			&i.WorkoutExerciseID,
			&i.WorkoutExerciseCreatedAt,
			&i.WorkoutExerciseUpdatedAt,
			&i.ExerciseID,
			&i.ExerciseName,
			&i.CategoryName,
			&i.BodyPartName,
			&i.SetID,
			&i.Weight,
			&i.DeductedWeight,
			&i.Duration,
			&i.Reps,
			&i.SetCreatedAt,
			&i.SetUpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkouts = `-- name: GetWorkouts :many
SELECT workouts.id, workouts.name, workouts.created_at, workouts.updated_at, workouts.start_time, workouts.end_time, users.id as user_id
FROM workouts
INNER JOIN users ON workouts.user_id = users.id
WHERE workouts.user_id = $1
`

type GetWorkoutsRow struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	StartTime sql.NullTime
	EndTime   sql.NullTime
	UserID    uuid.UUID
}

func (q *Queries) GetWorkouts(ctx context.Context, userID uuid.NullUUID) ([]GetWorkoutsRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkouts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkoutsRow
	for rows.Next() {
		var i GetWorkoutsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.StartTime,
			&i.EndTime,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
