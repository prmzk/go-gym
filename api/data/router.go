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
	Tx   *sql.DB
}

func NewApi(conn *sql.DB, authMiddleware func(http.Handler) http.Handler) (*dataApi, error) {
	db := dataStore.New(conn)

	dataApi := &dataApi{
		DB:   db,
		auth: authMiddleware,
		Tx:   conn,
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

	exerciseRouter.Post("/workout-user/{id}", dataApi.handlerUpsertExerciseUser)

	workoutsRouter := chi.NewRouter()
	workoutsRouter.Use(dataApi.auth)
	workoutsRouter.Get("/", dataApi.handlerGetWorkouts)
	workoutsRouter.Post("/", dataApi.handlerCreateWorkout)
	workoutsRouter.Get("/{id}", dataApi.handlerGetWorkoutByID)
	workoutsRouter.Delete("/{id}", dataApi.handleDeleteWorkout)

	workoutsRouter.Get("/previous-sets/{id}", dataApi.handlerGetPreviousWorkoutExerciseSets)

	r.Mount("/exercises", exerciseRouter)
	r.Mount("/workouts", workoutsRouter)

	return r
}
