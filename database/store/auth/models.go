// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Exercise struct {
	ID         uuid.UUID
	Name       string
	CategoryID uuid.NullUUID
	BodyPartID uuid.NullUUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ImageUrl   sql.NullString
	VideoUrl   sql.NullString
}

type ExerciseBodyPart struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExerciseCategory struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExerciseUser struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     uuid.UUID
	ExerciseID uuid.UUID
	Notes      sql.NullString
	RestTime   sql.NullInt32
}

type Set struct {
	ID                uuid.UUID
	WorkoutExerciseID uuid.NullUUID
	Weight            sql.NullFloat64
	DeductedWeight    sql.NullFloat64
	Duration          sql.NullInt32
	Reps              sql.NullInt32
	OrderNo           int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type SetTemplate struct {
	ID                 uuid.UUID
	TemplateExerciseID uuid.NullUUID
	Weight             sql.NullFloat64
	DeductedWeight     sql.NullFloat64
	Duration           sql.NullInt32
	Reps               sql.NullInt32
	OrderNo            int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Template struct {
	ID        uuid.UUID
	UserID    uuid.NullUUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TemplateExercise struct {
	ID         uuid.UUID
	TemplateID uuid.NullUUID
	ExerciseID uuid.NullUUID
	OrderNo    int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Token struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     uuid.NullUUID
	Expiration time.Time
	Type       string
}

type User struct {
	ID        uuid.UUID
	Name      sql.NullString
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Workout struct {
	ID        uuid.UUID
	UserID    uuid.NullUUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	StartTime sql.NullTime
	EndTime   sql.NullTime
}

type WorkoutExercise struct {
	ID         uuid.UUID
	WorkoutID  uuid.NullUUID
	ExerciseID uuid.NullUUID
	OrderNo    int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
