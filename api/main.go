package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/prmzk/go-base-prmzk/database"
)

type apiConfig struct {
	DB *database.Queries
}

func NewRouter(db *database.Queries) (*chi.Mux, error) {
	apiCfg := apiConfig{
		DB: db,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(corsConfig().Handler)

	r.Get("/", apiCfg.handlerGetUsers)
	r.Post("/", apiCfg.handlerCreateUser)

	// r.Get("/healthz", )

	return r, nil
}
