package data

import "errors"

var (
	ErrInvalidStartTime       = errors.New("invalid start time")
	ErrInvalidEndTime         = errors.New("invalid end time")
	ErrInvalidCreatedAt       = errors.New("invalid created at")
	ErrInvalidWorkoutExercise = errors.New("invalid workout exercise")
)
