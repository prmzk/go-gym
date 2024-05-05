package data

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/api/auth"
	"github.com/prmzk/go-base-prmzk/api/response"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
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
