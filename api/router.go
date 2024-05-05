package api

import (
	"database/sql"
	"log"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/prmzk/go-base-prmzk/api/auth"
	"github.com/prmzk/go-base-prmzk/api/data"
)

func NewRouter(url string) (*chi.Mux, error) {
	// apiCfg := &Api{
	// 	DB: db,
	// }

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(corsConfig().Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)

	// User/Auth
	authRouter, err := auth.NewApi(db)
	if err != nil {
		log.Fatal(err)
	}
	v1Router.Mount("/users", authRouter.Router())

	dataRouter, err := data.NewApi(db, authRouter.VerifyToken)
	if err != nil {
		log.Fatal(err)
	}
	v1Router.Mount("/gym", dataRouter.Router())

	r.Mount("/v1", v1Router)
	return r, nil
}
