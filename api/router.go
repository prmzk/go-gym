package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	"github.com/prmzk/go-base-prmzk/api/auth"
	"github.com/prmzk/go-base-prmzk/api/data"
	"github.com/prmzk/go-base-prmzk/api/response"
)

func NewRouter(url string) (*chi.Mux, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(corsConfig().Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, response.SuccessResponseOK(nil))
	})

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

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}
