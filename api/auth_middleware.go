package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/json"
)

type contextKey string

const userKey contextKey = "user"

func (apiCfg *apiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			json.RespondWithError(w, 401, "invalid authorization header")
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			json.RespondWithError(w, 401, "invalid authorization header")
			return
		}

		token := authParts[1]

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}
			// Return the secret key used for signing the token
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			json.RespondWithError(w, 401, "Invalid token")
			return
		}

		// Get the user ID from the token claims
		userID, ok := claims["sub"].(string)
		if !ok {
			json.RespondWithError(w, 401, "Invalid token")
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			json.RespondWithError(w, 400, "Invalid user ID")
			return
		}

		user, err := apiCfg.DB.GetUserById(r.Context(), userUUID)
		if err != nil {
			if err == sql.ErrNoRows {
				json.RespondWithError(w, 404, "User not found")
				return
			}
			json.RespondWithError(w, 500, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
