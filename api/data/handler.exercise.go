package data

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/api/auth"
	"github.com/prmzk/go-base-prmzk/api/response"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
	dataStore "github.com/prmzk/go-base-prmzk/database/store/data"
)

type Exercise struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageUrl  string    `json:"image_url,omitempty"`
	VideoUrl  string    `json:"video_url,omitempty"`
	Category  string    `json:"category,omitempty"`
	BodyPart  string    `json:"body_part,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	RestTime  int32     `json:"rest_time,omitempty"`
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BodyPart struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (dataApi *dataApi) handlerGetExercise(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	name := r.URL.Query().Get("name")
	category := r.URL.Query().Get("category")
	bodyPart := r.URL.Query().Get("body_part")

	if category == "all" {
		category = ""
	}

	if bodyPart == "all" {
		bodyPart = ""
	}

	exerciseRows, err := dataApi.DB.GetExercises(r.Context(), dataStore.GetExercisesParams{
		UserID:   user.ID,
		Name:     sql.NullString{String: name, Valid: name != ""},
		Category: sql.NullString{String: category, Valid: category != ""},
		BodyPart: sql.NullString{String: bodyPart, Valid: bodyPart != ""},
	})

	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	responseData := &struct {
		Exercises []Exercise `json:"exercises"`
	}{Exercises: []Exercise{}}
	for _, e := range exerciseRows {
		responseData.Exercises = append(responseData.Exercises, Exercise{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			ImageUrl:  e.ImageUrl.String,
			VideoUrl:  e.VideoUrl.String,
			Category:  e.Category.String,
			BodyPart:  e.BodyPart.String,
			Notes:     e.Notes.String,
			RestTime:  e.RestTime.Int32,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

func (dataApi *dataApi) handlerGetExcerciseById(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	exerciseID := chi.URLParam(r, "id")

	exerciseUUID, err := uuid.Parse(exerciseID)
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	exercise, err := dataApi.DB.GetExerciseById(r.Context(), dataStore.GetExerciseByIdParams{
		ID:     exerciseUUID,
		UserID: user.ID,
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	render.Render(w, r, response.SuccessResponseOK(&struct {
		Exercise Exercise `json:"exercise"`
	}{Exercise: Exercise{
		ID:        exercise.ID,
		Name:      exercise.Name,
		CreatedAt: exercise.CreatedAt,
		UpdatedAt: exercise.UpdatedAt,
		ImageUrl:  exercise.ImageUrl.String,
		VideoUrl:  exercise.VideoUrl.String,
		Category:  exercise.Category.String,
		BodyPart:  exercise.BodyPart.String,
		Notes:     exercise.Notes.String,
		RestTime:  exercise.RestTime.Int32,
	}}))
}

func (dataApi *dataApi) handlerGetCategories(w http.ResponseWriter, r *http.Request) {
	categoryRows, err := dataApi.DB.GetExerciseCategories(r.Context())
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	responseData := &struct {
		Categories []Category `json:"categories"`
	}{Categories: []Category{}}
	for _, e := range categoryRows {
		responseData.Categories = append(responseData.Categories, Category{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

func (dataApi *dataApi) handlerGetExcerciseByCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	categoryID := chi.URLParam(r, "id")

	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	exerciseRows, err := dataApi.DB.GetExerciseByCategory(r.Context(), dataStore.GetExerciseByCategoryParams{
		UserID:     user.ID,
		CategoryID: categoryUUID,
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	if len(exerciseRows) == 0 {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	responseData := &struct {
		ID        uuid.UUID  `json:"id"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		Exercises []Exercise `json:"exercises"`
	}{
		ID:        exerciseRows[0].CategoryID,
		Name:      exerciseRows[0].CategoryName,
		CreatedAt: exerciseRows[0].CategoryCreatedAt,
		UpdatedAt: exerciseRows[0].CategoryUpdatedAt,
		Exercises: []Exercise{},
	}

	if exerciseRows[0].ID.UUID == uuid.Nil {
		render.Render(w, r, response.SuccessResponseOK(responseData))
		return
	}

	for _, e := range exerciseRows {
		responseData.Exercises = append(responseData.Exercises, Exercise{
			ID:        e.ID.UUID,
			Name:      e.Name.String,
			CreatedAt: e.CreatedAt.Time,
			UpdatedAt: e.UpdatedAt.Time,
			ImageUrl:  e.ImageUrl.String,
			VideoUrl:  e.VideoUrl.String,
			BodyPart:  e.BodyPart.String,
			Notes:     e.Notes.String,
			RestTime:  e.RestTime.Int32,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

func (dataApi *dataApi) handlerGetBodyParts(w http.ResponseWriter, r *http.Request) {
	bodyPartRows, err := dataApi.DB.GetExerciseBodyParts(r.Context())
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	responseData := &struct {
		BodyParts []BodyPart `json:"body_parts"`
	}{BodyParts: []BodyPart{}}
	for _, e := range bodyPartRows {
		responseData.BodyParts = append(responseData.BodyParts, BodyPart{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

func (dataApi *dataApi) handlerGetExcerciseByBodyPart(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	bodyPartId := chi.URLParam(r, "id")

	bodyPartUUID, err := uuid.Parse(bodyPartId)
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	exerciseRows, err := dataApi.DB.GetExerciseByBodyPart(r.Context(), dataStore.GetExerciseByBodyPartParams{
		UserID:     user.ID,
		BodyPartID: bodyPartUUID,
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	if len(exerciseRows) == 0 {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	responseData := &struct {
		ID        uuid.UUID  `json:"id"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		Exercises []Exercise `json:"exercises"`
	}{
		ID:        exerciseRows[0].BodyPartID,
		Name:      exerciseRows[0].BodyPartName,
		CreatedAt: exerciseRows[0].BodyPartCreatedAt,
		UpdatedAt: exerciseRows[0].BodyPartUpdatedAt,
		Exercises: []Exercise{},
	}

	if exerciseRows[0].ID.UUID == uuid.Nil {
		render.Render(w, r, response.SuccessResponseOK(responseData))
		return
	}

	for _, e := range exerciseRows {
		responseData.Exercises = append(responseData.Exercises, Exercise{
			ID:        e.ID.UUID,
			Name:      e.Name.String,
			Category:  e.Category.String,
			CreatedAt: e.CreatedAt.Time,
			UpdatedAt: e.UpdatedAt.Time,
			ImageUrl:  e.ImageUrl.String,
			VideoUrl:  e.VideoUrl.String,
			Notes:     e.Notes.String,
			RestTime:  e.RestTime.Int32,
		})
	}

	render.Render(w, r, response.SuccessResponseOK(responseData))
}

type upsertExercuseUserRequest struct {
	Notes    *string `json:"notes"`
	RestTime *int    `json:"rest_time"`
}

func (body *upsertExercuseUserRequest) Bind(r *http.Request) error {
	return nil
}

func (dataApi *dataApi) handlerUpsertExerciseUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	exerciseID := chi.URLParam(r, "id")

	exerciseUUID, err := uuid.Parse(exerciseID)
	if err != nil {
		render.Render(w, r, response.ErrorResponseNotFound())
		return
	}

	upsertRequest := &upsertExercuseUserRequest{}
	if err := render.Bind(r, upsertRequest); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	var Notes sql.NullString
	var RestTime sql.NullInt32

	if upsertRequest.Notes == nil {
		Notes = sql.NullString{String: "", Valid: false}
	} else {
		Notes = sql.NullString{String: *upsertRequest.Notes, Valid: true}
	}

	if upsertRequest.RestTime == nil {
		RestTime = sql.NullInt32{Int32: 0, Valid: false}
	} else {
		RestTime = sql.NullInt32{Int32: int32(*upsertRequest.RestTime), Valid: true}
	}

	err = dataApi.DB.UpsertExerciseUser(r.Context(), dataStore.UpsertExerciseUserParams{
		UserID:     user.ID,
		ExerciseID: exerciseUUID,
		Notes:      Notes,
		RestTime:   RestTime,
	})

	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	render.Render(w, r, response.SuccessResponseOK(nil))

}
