package auth

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
	"github.com/prmzk/go-base-prmzk/api/response"
)

type contextKey string

const UserKey contextKey = "user"
const JwtSub = "sub"

func (authApi *authApi) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
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
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}

		// Get the user ID from the token claims
		userID, ok := claims[JwtSub].(string)
		if !ok {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			render.Render(w, r, response.ErrorResponseInternalServerError())
			return
		}

		user, err := authApi.DB.GetUserById(r.Context(), userUUID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
				return
			}
			render.Render(w, r, response.ErrorResponseInternalServerError())
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
