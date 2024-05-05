package auth

import (
	"database/sql"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
)

type authApi struct {
	DB *authStore.Queries
}

func NewApi(conn *sql.DB) (*authApi, error) {
	db := authStore.New(conn)

	authApi := &authApi{
		DB: db,
	}

	authApi.choresTicker()

	return authApi, nil
}

func (authApi *authApi) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", authApi.handlerCreateUser)
	r.Post("/login", authApi.handlerLoginUser)
	r.Get("/login/callback", authApi.handlerValidateToken)
	r.Get("/refresh", authApi.handlerRefreshToken)
	r.With(authApi.VerifyToken).Get("/me", authApi.handlerMe)

	return r
}
