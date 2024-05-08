package data

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	dataStore "github.com/prmzk/go-base-prmzk/database/store/data"
)

type dataApi struct {
	DB   *dataStore.Queries
	auth func(http.Handler) http.Handler
}

func NewApi(conn *sql.DB, authMiddleware func(http.Handler) http.Handler) (*dataApi, error) {
	db := dataStore.New(conn)

	dataApi := &dataApi{
		DB:   db,
		auth: authMiddleware,
	}
	return dataApi, nil
}

func (dataApi *dataApi) Router() *chi.Mux {
	r := chi.NewRouter()

	exerciseRouter := chi.NewRouter()
	exerciseRouter.Use(dataApi.auth)
	exerciseRouter.Get("/", dataApi.handlerGetExercise)
	exerciseRouter.Get("/{id}", dataApi.handlerGetExcerciseById)
	exerciseRouter.Get("/categories*", dataApi.handlerGetCategories)
	exerciseRouter.Get("/categories/{id}", dataApi.handlerGetExcerciseByCategory)
	exerciseRouter.Get("/bodyparts*", dataApi.handlerGetBodyParts)
	exerciseRouter.Get("/bodyparts/{id}", dataApi.handlerGetExcerciseByBodyPart)

	workoutsRouter := chi.NewRouter()
	workoutsRouter.Use(dataApi.auth)
	workoutsRouter.Get("/", dataApi.handlerGetWorkouts)

	r.Mount("/exercises", exerciseRouter)
	r.Mount("/workouts", workoutsRouter)

	return r
}
