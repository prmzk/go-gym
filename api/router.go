package api

import (
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/prmzk/go-base-prmzk/database"
)

type apiConfig struct {
	DB *database.Queries
}

func NewRouter(db *database.Queries) (*chi.Mux, error) {
	apiCfg := &apiConfig{
		DB: db,
	}

	r := chi.NewRouter()

	// r.Use(middleware.Logger)
	r.Use(corsConfig().Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)

	userRouter := chi.NewRouter()

	userRouter.Post("/register", apiCfg.handlerCreateUser)
	userRouter.Post("/login", apiCfg.handlerLoginUser)
	userRouter.Get("/login/callback", apiCfg.handlerValidateToken)

	userRouter.With(apiCfg.AuthMiddleware).Get("/me", apiCfg.handlerMe)

	v1Router.Mount("/users", userRouter)
	r.Mount("/v1", v1Router)

	return r, nil
}
