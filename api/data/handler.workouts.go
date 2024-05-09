package data

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/api/auth"
	"github.com/prmzk/go-base-prmzk/api/response"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
	"github.com/prmzk/go-base-prmzk/database/store/data"
)

type Workout struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (dataApi *dataApi) handlerGetWorkouts(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	workoutRows, err := dataApi.DB.GetWorkouts(r.Context(), uuid.NullUUID{UUID: user.ID, Valid: true})

	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	responseData := &struct {
		ID       uuid.UUID `json:"id"`
		Workouts []Workout `json:"workouts"`
	}{
		ID:       user.ID,
		Workouts: []Workout{},
	}
	for _, e := range workoutRows {
		responseData.Workouts = append(responseData.Workouts, Workout{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			StartTime: e.StartTime.Time,
			EndTime:   e.EndTime.Time,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

type GetWorkoutDetail struct {
	ID               uuid.UUID     `json:"id"`
	Name             string        `json:"name"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
	WorkoutExercises []GetExercise `json:"workout_exercises"`
}

type GetExercise struct {
	WorkoutExerciseID        uuid.UUID `json:"workout_exercise_id"`
	WorkoutExerciseCreatedAt time.Time `json:"workout_exercise_created_at"`
	WorkoutExerciseUpdatedAt time.Time `json:"workout_exercise_updated_at"`
	ExerciseID               uuid.UUID `json:"exercise_id"`
	ExerciseName             string    `json:"exercise_name"`
	CategoryName             string    `json:"category_name"`
	BodyPartName             string    `json:"body_part_name"`
	Sets                     []Set     `json:"sets"`
}

type Set struct {
	SetID          uuid.UUID `json:"set_id"`
	Weight         float64   `json:"weight,omitempty"`
	DeductedWeight float64   `json:"deducted_weight,omitempty"`
	Duration       int32     `json:"duration,omitempty"`
	Reps           int32     `json:"reps,omitempty"`
	SetCreatedAt   time.Time `json:"set_created_at"`
	SetUpdatedAt   time.Time `json:"set_updated_at"`
}

func (dataApi *dataApi) handlerGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	workoutRow, err := dataApi.DB.GetWorkoutById(r.Context(), data.GetWorkoutByIdParams{
		ID:     workoutID,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}
	if len(workoutRow) == 0 {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	parsedWorkout := GetWorkoutDetail{
		ID:               workoutRow[0].ID,
		Name:             workoutRow[0].Name,
		CreatedAt:        workoutRow[0].CreatedAt,
		UpdatedAt:        workoutRow[0].UpdatedAt,
		StartTime:        workoutRow[0].StartTime.Time,
		EndTime:          workoutRow[0].EndTime.Time,
		WorkoutExercises: []GetExercise{},
	}
	// Create a map for quick lookup of exercises
	exerciseIndex := make(map[uuid.UUID]int)

	for _, e := range workoutRow {
		if i, exists := exerciseIndex[e.ExerciseID]; exists {
			// Exercise exists, append the set
			parsedWorkout.WorkoutExercises[i].Sets = append(parsedWorkout.WorkoutExercises[i].Sets, Set{
				SetID:          e.SetID,
				Weight:         e.Weight.Float64,
				DeductedWeight: e.DeductedWeight.Float64,
				Duration:       e.Duration.Int32,
				Reps:           e.Reps.Int32,
				SetCreatedAt:   e.SetCreatedAt,
				SetUpdatedAt:   e.SetUpdatedAt,
			})
		} else {
			// Exercise doesn't exist, create a new one
			newExercise := GetExercise{
				WorkoutExerciseID:        e.WorkoutExerciseID,
				WorkoutExerciseCreatedAt: e.WorkoutExerciseCreatedAt,
				WorkoutExerciseUpdatedAt: e.WorkoutExerciseUpdatedAt,
				ExerciseID:               e.ExerciseID,
				ExerciseName:             e.ExerciseName,
				CategoryName:             e.CategoryName,
				BodyPartName:             e.BodyPartName,
				Sets: []Set{
					{
						SetID:          e.SetID,
						Weight:         e.Weight.Float64,
						DeductedWeight: e.DeductedWeight.Float64,
						Duration:       e.Duration.Int32,
						Reps:           e.Reps.Int32,
						SetCreatedAt:   e.SetCreatedAt,
						SetUpdatedAt:   e.SetUpdatedAt,
					},
				},
			}
			parsedWorkout.WorkoutExercises = append(parsedWorkout.WorkoutExercises, newExercise)
			// Store the index of the new exercise in the map
			exerciseIndex[e.ExerciseID] = len(parsedWorkout.WorkoutExercises) - 1
		}
	}

	render.Render(w, r, response.SuccessResponseOK(parsedWorkout))

}

type workoutRequest struct {
	Workout          workoutDetailRequest     `json:"workout"`
	WorkoutExercises []workoutExerciseRequest `json:"workout_exercises"`
	Sets             []setRequest             `json:"sets"`
}

type workoutDetailRequest struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
}

type workoutExerciseRequest struct {
	ID         uuid.UUID `json:"id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type setRequest struct {
	WorkoutExerciseID uuid.UUID `json:"workout_exercise_id"`
	ID                uuid.UUID `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	Weight            float32   `json:"weight"`
	DeductedWeight    float32   `json:"deducted_weight"`
	Duration          int32     `json:"duration"`
	Reps              int32     `json:"reps"`
}

func (body *workoutRequest) Bind(r *http.Request) error {
	if body.Workout.Name == "" {
		return fmt.Errorf("name is required")
	}

	if body.Workout.EndTime.Before(body.Workout.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}

func (dataApi *dataApi) handlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	tx, err := dataApi.Tx.Begin()
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}
	defer tx.Rollback()
	qtx := dataApi.DB.WithTx(tx)

	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	workoutDetailRequest := &workoutRequest{}
	if err := render.Bind(r, workoutDetailRequest); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	workoutRow, err := qtx.CreateWorkout(r.Context(), data.CreateWorkoutParams{
		Name:      workoutDetailRequest.Workout.Name,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		StartTime: sql.NullTime{Time: workoutDetailRequest.Workout.StartTime.UTC(), Valid: true},
		EndTime:   sql.NullTime{Time: workoutDetailRequest.Workout.EndTime.UTC(), Valid: true},
		CreatedAt: workoutDetailRequest.Workout.CreatedAt.UTC(),
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	createWorkoutExerciseParams := data.CreateWorkoutExerciseParams{
		IDArray:         []uuid.UUID{},
		WorkoutIDArray:  []uuid.UUID{},
		ExerciseIDArray: []uuid.UUID{},
		CreatedAtArray:  []time.Time{},
	}

	for _, e := range workoutDetailRequest.WorkoutExercises {
		createWorkoutExerciseParams.IDArray = append(createWorkoutExerciseParams.IDArray, e.ID)
		createWorkoutExerciseParams.WorkoutIDArray = append(createWorkoutExerciseParams.WorkoutIDArray, workoutRow.ID)
		createWorkoutExerciseParams.ExerciseIDArray = append(createWorkoutExerciseParams.ExerciseIDArray, e.ExerciseID)
		createWorkoutExerciseParams.CreatedAtArray = append(createWorkoutExerciseParams.CreatedAtArray, e.CreatedAt.UTC())
	}

	_, err = qtx.CreateWorkoutExercise(r.Context(), createWorkoutExerciseParams)
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(ErrInvalidWorkoutExercise))
		return
	}

	createSetsParams := data.CreateSetsParams{
		IDArray:                []uuid.UUID{},
		WorkoutExerciseIDArray: []uuid.UUID{},
		WeightArray:            []float32{},
		DeductedWeightArray:    []float32{},
		DurationArray:          []int32{},
		RepsArray:              []int32{},
		CreatedAtArray:         []time.Time{},
	}

	for _, e := range workoutDetailRequest.Sets {
		createSetsParams.IDArray = append(createSetsParams.IDArray, e.ID)
		createSetsParams.WorkoutExerciseIDArray = append(createSetsParams.WorkoutExerciseIDArray, e.WorkoutExerciseID)
		createSetsParams.WeightArray = append(createSetsParams.WeightArray, e.Weight)
		createSetsParams.DeductedWeightArray = append(createSetsParams.DeductedWeightArray, e.DeductedWeight)
		createSetsParams.DurationArray = append(createSetsParams.DurationArray, e.Duration)
		createSetsParams.RepsArray = append(createSetsParams.RepsArray, e.Reps)
		createSetsParams.CreatedAtArray = append(createSetsParams.CreatedAtArray, e.CreatedAt.UTC())
	}

	_, err = qtx.CreateSets(r.Context(), createSetsParams)
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	if err = tx.Commit(); err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	render.Render(w, r, response.SuccessResponseOK(nil))
}
