package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userKey contextKey = "user"

func (apiCfg *apiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
			return
		}

		token := authParts[1]

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
			return
		}

		// Get the user ID from the token claims
		userID, ok := claims["sub"].(string)
		if !ok {
			render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			render.Render(w, r, ErrInternalServerError)
			return
		}

		user, err := apiCfg.DB.GetUserById(r.Context(), userUUID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
				return
			}
			render.Render(w, r, ErrInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
