package data

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/api/auth"
	"github.com/prmzk/go-base-prmzk/api/response"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
	"github.com/prmzk/go-base-prmzk/database/store/data"
)

type Template struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (dataApi *dataApi) handleGetTemplate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	templateRows, err := dataApi.DB.GetTemplates(r.Context(), uuid.NullUUID{UUID: user.ID, Valid: true})

	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	responseData := &struct {
		ID       uuid.UUID  `json:"id"`
		Workouts []Template `json:"workouts"`
	}{
		ID:       user.ID,
		Workouts: []Template{},
	}
	for _, e := range templateRows {
		responseData.Workouts = append(responseData.Workouts, Template{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

type GetTemplateDetail struct {
	ID               uuid.UUID             `json:"id"`
	Name             string                `json:"name"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	WorkoutExercises []GetExerciseTemplate `json:"workout_exercises"`
}

type GetExerciseTemplate struct {
	TemplateExerciseId          uuid.UUID        `json:"workout_exercise_id"`
	TemplateExerciseIdCreatedAt time.Time        `json:"workout_exercise_created_at"`
	TemplateExerciseIdUpdatedAt time.Time        `json:"workout_exercise_updated_at"`
	ExerciseID                  uuid.UUID        `json:"exercise_id"`
	ExerciseName                string           `json:"exercise_name"`
	CategoryName                string           `json:"category_name"`
	BodyPartName                string           `json:"body_part_name"`
	Sets                        []GetSetTemplate `json:"sets"`
	OrderNo                     int32            `json:"order_no"`
}

type GetSetTemplate struct {
	SetID          uuid.UUID `json:"set_id"`
	Weight         float64   `json:"weight,omitempty"`
	DeductedWeight float64   `json:"deducted_weight,omitempty"`
	Duration       int32     `json:"duration,omitempty"`
	Reps           int32     `json:"reps,omitempty"`
	SetCreatedAt   time.Time `json:"set_created_at"`
	SetUpdatedAt   time.Time `json:"set_updated_at"`
	OrderNo        int32     `json:"order_no"`
}

func (dataApi *dataApi) handleGetTemplateById(w http.ResponseWriter, r *http.Request) {
	templateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	templateRow, err := dataApi.DB.GetTemplateById(r.Context(), data.GetTemplateByIdParams{
		ID:     templateID,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrorResponseNotFound())
			return
		}
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	if len(templateRow) == 0 {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	parsedTemplate := GetTemplateDetail{
		ID:               templateRow[0].ID,
		Name:             templateRow[0].Name,
		CreatedAt:        templateRow[0].CreatedAt,
		UpdatedAt:        templateRow[0].UpdatedAt,
		WorkoutExercises: []GetExerciseTemplate{},
	}
	// Create a map for quick lookup of exercises
	exerciseIndex := make(map[uuid.UUID]int)

	for _, e := range templateRow {
		if i, exists := exerciseIndex[e.ExerciseID]; exists {
			// Exercise exists, append the set
			parsedTemplate.WorkoutExercises[i].Sets = append(parsedTemplate.WorkoutExercises[i].Sets, GetSetTemplate{
				SetID:          e.SetID.UUID,
				Weight:         e.Weight.Float64,
				DeductedWeight: e.DeductedWeight.Float64,
				Duration:       e.Duration.Int32,
				Reps:           e.Reps.Int32,
				SetCreatedAt:   e.SetCreatedAt.Time,
				SetUpdatedAt:   e.SetUpdatedAt.Time,
				OrderNo:        e.SetTemplatesOrderNo.Int32,
			})
		} else {
			// Exercise doesn't exist, create a new one
			newExercise := GetExerciseTemplate{
				TemplateExerciseId:          e.TemplateExerciseID,
				TemplateExerciseIdCreatedAt: e.WorkoutExerciseCreatedAt,
				TemplateExerciseIdUpdatedAt: e.WorkoutExerciseUpdatedAt,
				ExerciseID:                  e.ExerciseID,
				ExerciseName:                e.ExerciseName,
				CategoryName:                e.CategoryName,
				BodyPartName:                e.BodyPartName,
				OrderNo:                     e.TemplateExerciseOrderNo,
				Sets:                        []GetSetTemplate{},
			}
			if e.SetID.Valid {
				newExercise.Sets = append(newExercise.Sets, GetSetTemplate{
					SetID:          e.SetID.UUID,
					Weight:         e.Weight.Float64,
					DeductedWeight: e.DeductedWeight.Float64,
					Duration:       e.Duration.Int32,
					Reps:           e.Reps.Int32,
					SetCreatedAt:   e.SetCreatedAt.Time,
					SetUpdatedAt:   e.SetUpdatedAt.Time,
					OrderNo:        e.SetTemplatesOrderNo.Int32,
				})
			}
			parsedTemplate.WorkoutExercises = append(parsedTemplate.WorkoutExercises, newExercise)
			// Store the index of the new exercise in the map
			exerciseIndex[e.ExerciseID] = len(parsedTemplate.WorkoutExercises) - 1
		}
	}

	// Sort the workout exercises based on orderNo
	sort.Slice(parsedTemplate.WorkoutExercises, func(i, j int) bool {
		return parsedTemplate.WorkoutExercises[i].OrderNo < parsedTemplate.WorkoutExercises[j].OrderNo
	})

	// Sort the sets inside each workout exercise based on orderNo
	for _, exercise := range parsedTemplate.WorkoutExercises {
		sort.Slice(exercise.Sets, func(i, j int) bool {
			return exercise.Sets[i].OrderNo < exercise.Sets[j].OrderNo
		})
	}

	render.Render(w, r, response.SuccessResponseOK(parsedTemplate))
}

type templateRequest struct {
	Workout          templateDetailRequest            `json:"workout"`
	WorkoutExercises []templateWorkoutExerciseRequest `json:"workout_exercises"`
	Sets             []setTemplateRequest             `json:"sets"`
}

type templateDetailRequest struct {
	Name string `json:"name"`
}

type templateWorkoutExerciseRequest struct {
	ID         uuid.UUID `json:"id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
	OrderNo    int32     `json:"order_no"`
}

type setTemplateRequest struct {
	ID                uuid.UUID `json:"id"`
	WorkoutExerciseID uuid.UUID `json:"template_exercise_id"`
	Weight            float32   `json:"weight"`
	DeductedWeight    float32   `json:"deducted_weight"`
	Duration          int32     `json:"duration"`
	Reps              int32     `json:"reps"`
	OrderNo           int32     `json:"order_no"`
}

func (body *templateRequest) Bind(r *http.Request) error {
	if body.Workout.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (dataApi *dataApi) handleCreateTemplate(w http.ResponseWriter, r *http.Request) {
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

	templateDetailRequest := &templateRequest{}
	if err := render.Bind(r, templateDetailRequest); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	templateRow, err := qtx.CreateTemplate(r.Context(), data.CreateTemplateParams{
		Name:   templateDetailRequest.Workout.Name,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	createTemplateExerciseParams := data.CreateTemplateExercisesParams{
		IDArray:         []uuid.UUID{},
		TemplateID:      []uuid.UUID{},
		ExerciseIDArray: []uuid.UUID{},
		OrderNoArray:    []int32{},
	}

	for _, e := range templateDetailRequest.WorkoutExercises {
		createTemplateExerciseParams.IDArray = append(createTemplateExerciseParams.IDArray, e.ID)
		createTemplateExerciseParams.TemplateID = append(createTemplateExerciseParams.TemplateID, templateRow.ID)
		createTemplateExerciseParams.ExerciseIDArray = append(createTemplateExerciseParams.ExerciseIDArray, e.ExerciseID)
		createTemplateExerciseParams.OrderNoArray = append(createTemplateExerciseParams.OrderNoArray, e.OrderNo)
	}

	_, err = qtx.CreateTemplateExercises(r.Context(), createTemplateExerciseParams)
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	createSetsTemplateParams := data.CreateSetTemplatesParams{
		IDArray:             []uuid.UUID{},
		TemplateExerciseID:  []uuid.UUID{},
		WeightArray:         []float32{},
		DeductedWeightArray: []float32{},
		DurationArray:       []int32{},
		RepsArray:           []int32{},
		OrderNoArray:        []int32{},
	}

	for _, e := range templateDetailRequest.Sets {
		createSetsTemplateParams.IDArray = append(createSetsTemplateParams.IDArray, e.ID)
		createSetsTemplateParams.TemplateExerciseID = append(createSetsTemplateParams.TemplateExerciseID, e.WorkoutExerciseID)
		createSetsTemplateParams.WeightArray = append(createSetsTemplateParams.WeightArray, e.Weight)
		createSetsTemplateParams.DeductedWeightArray = append(createSetsTemplateParams.DeductedWeightArray, e.DeductedWeight)
		createSetsTemplateParams.DurationArray = append(createSetsTemplateParams.DurationArray, e.Duration)
		createSetsTemplateParams.RepsArray = append(createSetsTemplateParams.RepsArray, e.Reps)
		createSetsTemplateParams.OrderNoArray = append(createSetsTemplateParams.OrderNoArray, e.OrderNo)
	}

	_, err = qtx.CreateSetTemplates(r.Context(), createSetsTemplateParams)
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

func (dataApi *dataApi) handleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)
	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	err = dataApi.DB.DeleteTemplate(r.Context(), data.DeleteTemplateParams{
		TemplateID: templateId,
		UserID:     uuid.NullUUID{UUID: user.ID, Valid: true},
	})

	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	render.Render(w, r, response.SuccessResponseOK(nil))
}

func (dataApi *dataApi) handleChangeSetTemplate(w http.ResponseWriter, r *http.Request) {
	templateId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)
	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	tx, err := dataApi.Tx.Begin()
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}
	defer tx.Rollback()
	qtx := dataApi.DB.WithTx(tx)

	workoutDetailRequest := &workoutRequest{}
	if err := render.Bind(r, workoutDetailRequest); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	err = qtx.DeleteTemplateExercise(r.Context(), data.DeleteTemplateExerciseParams{
		TemplateID: templateId,
		UserID:     uuid.NullUUID{UUID: user.ID, Valid: true},
	})

	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	createTemplateExerciseParams := data.CreateTemplateExercisesParams{
		IDArray:         []uuid.UUID{},
		TemplateID:      []uuid.UUID{},
		ExerciseIDArray: []uuid.UUID{},
		OrderNoArray:    []int32{},
	}

	for _, e := range workoutDetailRequest.WorkoutExercises {
		createTemplateExerciseParams.IDArray = append(createTemplateExerciseParams.IDArray, e.ID)
		createTemplateExerciseParams.TemplateID = append(createTemplateExerciseParams.TemplateID, templateId)
		createTemplateExerciseParams.ExerciseIDArray = append(createTemplateExerciseParams.ExerciseIDArray, e.ExerciseID)
		createTemplateExerciseParams.OrderNoArray = append(createTemplateExerciseParams.OrderNoArray, e.OrderNo)
	}

	_, err = qtx.CreateTemplateExercises(r.Context(), createTemplateExerciseParams)
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	createSetsTemplateParams := data.CreateSetTemplatesParams{
		IDArray:             []uuid.UUID{},
		TemplateExerciseID:  []uuid.UUID{},
		WeightArray:         []float32{},
		DeductedWeightArray: []float32{},
		DurationArray:       []int32{},
		RepsArray:           []int32{},
		OrderNoArray:        []int32{},
	}

	for _, e := range workoutDetailRequest.Sets {
		createSetsTemplateParams.IDArray = append(createSetsTemplateParams.IDArray, e.ID)
		createSetsTemplateParams.TemplateExerciseID = append(createSetsTemplateParams.TemplateExerciseID, e.WorkoutExerciseID)
		createSetsTemplateParams.WeightArray = append(createSetsTemplateParams.WeightArray, e.Weight)
		createSetsTemplateParams.DeductedWeightArray = append(createSetsTemplateParams.DeductedWeightArray, e.DeductedWeight)
		createSetsTemplateParams.DurationArray = append(createSetsTemplateParams.DurationArray, e.Duration)
		createSetsTemplateParams.RepsArray = append(createSetsTemplateParams.RepsArray, e.Reps)
		createSetsTemplateParams.OrderNoArray = append(createSetsTemplateParams.OrderNoArray, e.OrderNo)
	}

	_, err = qtx.CreateSetTemplates(r.Context(), createSetsTemplateParams)
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

func (dataApi *dataApi) handleChangeSetTemplateValueOnly(w http.ResponseWriter, r *http.Request) {
	templateId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)
	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	tx, err := dataApi.Tx.Begin()
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}
	defer tx.Rollback()
	qtx := dataApi.DB.WithTx(tx)

	workoutDetailRequest := &workoutRequest{}
	if err := render.Bind(r, workoutDetailRequest); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	// extract all unique exercise id from workoutDetailRequest.sets
	uniqueExerciseIds := []uuid.UUID{}
	for _, e := range workoutDetailRequest.WorkoutExercises {
		// check if exercise id already exists in uniqueExerciseIds
		exists := false
		for _, id := range uniqueExerciseIds {
			if id == e.ExerciseID {
				exists = true
				break
			}
		}
		if !exists {
			uniqueExerciseIds = append(uniqueExerciseIds, e.ExerciseID)
		}
	}

	affectedTemplateExercises, err := qtx.GetTemplateExerciseByExerciseId(r.Context(), uniqueExerciseIds)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	err = qtx.DeleteSetTemplate(r.Context(), data.DeleteSetTemplateParams{
		TemplateID: templateId,
		UserID:     uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	createSetsTemplateParams := data.CreateSetTemplatesParams{
		IDArray:             []uuid.UUID{},
		TemplateExerciseID:  []uuid.UUID{},
		WeightArray:         []float32{},
		DeductedWeightArray: []float32{},
		DurationArray:       []int32{},
		RepsArray:           []int32{},
		OrderNoArray:        []int32{},
	}

	for _, eWe := range workoutDetailRequest.WorkoutExercises {
		for _, eSet := range workoutDetailRequest.Sets {
			for _, eAffected := range affectedTemplateExercises {
				if eWe.ID == eSet.WorkoutExerciseID && eWe.ExerciseID == eAffected.ExerciseID.UUID {
					createSetsTemplateParams.IDArray = append(createSetsTemplateParams.IDArray, uuid.New())
					createSetsTemplateParams.TemplateExerciseID = append(createSetsTemplateParams.TemplateExerciseID, eAffected.ID)
					createSetsTemplateParams.WeightArray = append(createSetsTemplateParams.WeightArray, eSet.Weight)
					createSetsTemplateParams.DeductedWeightArray = append(createSetsTemplateParams.DeductedWeightArray, eSet.DeductedWeight)
					createSetsTemplateParams.DurationArray = append(createSetsTemplateParams.DurationArray, eSet.Duration)
					createSetsTemplateParams.RepsArray = append(createSetsTemplateParams.RepsArray, eSet.Reps)
					createSetsTemplateParams.OrderNoArray = append(createSetsTemplateParams.OrderNoArray, eSet.OrderNo)
				}
			}
		}
	}

	_, err = qtx.CreateSetTemplates(r.Context(), createSetsTemplateParams)
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
