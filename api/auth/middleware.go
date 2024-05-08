package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/api/response"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
)

type contextKey string

const UserKey contextKey = "user"

func (authApi *authApi) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
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
		// Parse the token
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			render.Render(w, r, response.ErrorResponseUnauthorized(err))
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}

		// Get user from token
		tokenUUID, err := uuid.Parse(jti)
		if err != nil {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}

		// Chek token validity
		tokenDB, err := authApi.DB.GetUserToken(r.Context(), tokenUUID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, response.ErrorResponseUnauthorized(ErrAccessTokenNotFound))
				return
			}
			render.Render(w, r, response.ErrorResponseInternalServerError())
			return
		}
		tokenIsNotValid := tokenDB.Type != "access" || tokenDB.Expiration.Before(time.Now())
		if tokenIsNotValid {
			_, err := authApi.DB.ClearUserToken(r.Context(), authStore.ClearUserTokenParams{
				TokenID: tokenUUID,
				Type:    "access",
			})
			if err != nil {
				render.Render(w, r, response.ErrorResponseInternalServerError())
				return
			}

			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}

		// Send user data to the next handler
		ctx := context.WithValue(r.Context(), UserKey, authStore.User{
			ID:        tokenDB.UserID.UUID,
			Name:      tokenDB.UserName,
			Email:     tokenDB.UserEmail,
			CreatedAt: tokenDB.UserCreatedAt,
			UpdatedAt: tokenDB.UserUpdatedAt,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
