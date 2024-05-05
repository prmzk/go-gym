// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: exercises.sql

package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getExerciseBodyParts = `-- name: GetExerciseBodyParts :many
SELECT id, name, created_at, updated_at FROM exercise_body_parts
`

func (q *Queries) GetExerciseBodyParts(ctx context.Context) ([]ExerciseBodyPart, error) {
	rows, err := q.db.QueryContext(ctx, getExerciseBodyParts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExerciseBodyPart
	for rows.Next() {
		var i ExerciseBodyPart
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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

const getExerciseByBodyPart = `-- name: GetExerciseByBodyPart :many
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, category.name as category, body_part.name as body_part_name, body_part.id as body_part_id, body_part.created_at as body_part_created_at, body_part.updated_at as body_part_updated_at
FROM exercises
RIGHT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
LEFT JOIN exercise_categories as category ON exercises.category_id = category.id
WHERE body_part.id = $1
`

type GetExerciseByBodyPartRow struct {
	ID                uuid.NullUUID
	Name              sql.NullString
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
	ImageUrl          sql.NullString
	VideoUrl          sql.NullString
	Category          sql.NullString
	BodyPartName      string
	BodyPartID        uuid.UUID
	BodyPartCreatedAt time.Time
	BodyPartUpdatedAt time.Time
}

func (q *Queries) GetExerciseByBodyPart(ctx context.Context, bodyPartID uuid.UUID) ([]GetExerciseByBodyPartRow, error) {
	rows, err := q.db.QueryContext(ctx, getExerciseByBodyPart, bodyPartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExerciseByBodyPartRow
	for rows.Next() {
		var i GetExerciseByBodyPartRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageUrl,
			&i.VideoUrl,
			&i.Category,
			&i.BodyPartName,
			&i.BodyPartID,
			&i.BodyPartCreatedAt,
			&i.BodyPartUpdatedAt,
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

const getExerciseByCategory = `-- name: GetExerciseByCategory :many
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, body_part.name as body_part, category.name as category_name, category.id as category_id, category.created_at as category_created_at, category.updated_at as category_updated_at
FROM exercises
RIGHT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
WHERE category.id = $1
`

type GetExerciseByCategoryRow struct {
	ID                uuid.NullUUID
	Name              sql.NullString
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
	ImageUrl          sql.NullString
	VideoUrl          sql.NullString
	BodyPart          sql.NullString
	CategoryName      string
	CategoryID        uuid.UUID
	CategoryCreatedAt time.Time
	CategoryUpdatedAt time.Time
}

func (q *Queries) GetExerciseByCategory(ctx context.Context, categoryID uuid.UUID) ([]GetExerciseByCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getExerciseByCategory, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExerciseByCategoryRow
	for rows.Next() {
		var i GetExerciseByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageUrl,
			&i.VideoUrl,
			&i.BodyPart,
			&i.CategoryName,
			&i.CategoryID,
			&i.CategoryCreatedAt,
			&i.CategoryUpdatedAt,
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

const getExerciseById = `-- name: GetExerciseById :one
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, category.name as category, body_part.name as body_part FROM exercises
LEFT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
WHERE exercises.id = $1
`

type GetExerciseByIdRow struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ImageUrl  sql.NullString
	VideoUrl  sql.NullString
	Category  sql.NullString
	BodyPart  sql.NullString
}

func (q *Queries) GetExerciseById(ctx context.Context, id uuid.UUID) (GetExerciseByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getExerciseById, id)
	var i GetExerciseByIdRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ImageUrl,
		&i.VideoUrl,
		&i.Category,
		&i.BodyPart,
	)
	return i, err
}

const getExerciseCategories = `-- name: GetExerciseCategories :many
SELECT id, name, created_at, updated_at FROM exercise_categories
`

func (q *Queries) GetExerciseCategories(ctx context.Context) ([]ExerciseCategory, error) {
	rows, err := q.db.QueryContext(ctx, getExerciseCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExerciseCategory
	for rows.Next() {
		var i ExerciseCategory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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

const getExercises = `-- name: GetExercises :many
SELECT exercises.id, exercises.name, exercises.created_at, exercises.updated_at, exercises.image_url, exercises.video_url, category.name as category, body_part.name as body_part 
FROM exercises
LEFT JOIN exercise_categories as category ON exercises.category_id = category.id
LEFT JOIN exercise_body_parts as body_part ON exercises.body_part_id = body_part.id
WHERE (exercises.name ILIKE '%' || $1 || '%' OR $1 IS NULL)
AND (category.name = $2 OR $2 IS NULL)
AND (body_part.name = $3 OR $3 IS NULL)
`

type GetExercisesParams struct {
	Name     sql.NullString
	Category sql.NullString
	BodyPart sql.NullString
}

type GetExercisesRow struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ImageUrl  sql.NullString
	VideoUrl  sql.NullString
	Category  sql.NullString
	BodyPart  sql.NullString
}

func (q *Queries) GetExercises(ctx context.Context, arg GetExercisesParams) ([]GetExercisesRow, error) {
	rows, err := q.db.QueryContext(ctx, getExercises, arg.Name, arg.Category, arg.BodyPart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExercisesRow
	for rows.Next() {
		var i GetExercisesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageUrl,
			&i.VideoUrl,
			&i.Category,
			&i.BodyPart,
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
