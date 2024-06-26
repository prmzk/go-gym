package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/api/response"
	authStore "github.com/prmzk/go-base-prmzk/database/store/auth"
)

var (
	refreshTokenExp = time.Hour * 24 * 7
	accessTokenExp  = time.Hour * 3
	loginTokenExp   = time.Hour * 1
)

func isValidEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	if !match {
		return ErrInvalidEmail
	}
	return nil
}

type registerRequest struct {
	Email string
	Name  string
}

func (body *registerRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	if body.Name == "" {
		return errors.New("name is required")
	}

	return isValidEmail(body.Email)
}

func (authApi *authApi) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Get Email from body
	body := &registerRequest{}
	if err := render.Bind(r, body); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(ErrInvalidEmailOrName))
		return
	}

	// Check if email is already in use
	_, err := authApi.DB.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			render.Render(w, r, response.ErrorResponseInternalServerError())
			return
		}
	} else {
		render.Render(w, r, response.ErrorResponseBadRequest(ErrDuplicateEmail))
		return
	}

	// Create user
	_, err = authApi.DB.CreateUser(r.Context(), authStore.CreateUserParams{
		ID:    uuid.New(),
		Email: body.Email,
		Name:  sql.NullString{String: body.Name, Valid: true},
	})
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	// Return success
	render.Render(w, r, response.SuccessResponseCreated(nil))
}

type loginRequest struct {
	Email string
}

func (body *loginRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	return isValidEmail(body.Email)
}

func (authApi *authApi) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	// Get Email from body
	body := &loginRequest{}
	if err := render.Bind(r, body); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(ErrInvalidEmail))
		return
	}

	// Check if email exists
	user, err := authApi.DB.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrorResponseBadRequest(ErrEmailNotFound))
			return
		}
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	// Set token jwtId and token expiration and save to db
	tokenClaims := &authStore.SetUserTokenParams{
		UserID:     uuid.NullUUID{UUID: user.ID, Valid: true},
		Type:       "login",
		Expiration: time.Now().Add(loginTokenExp),
	}

	_, err = authApi.DB.ClearAllTokenUser(r.Context(), tokenClaims.UserID)
	if err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(err))
		return
	}

	tokenDB, err := authApi.DB.SetUserToken(r.Context(), *tokenClaims)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	// Create token to be sent
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": tokenDB.Expiration.Unix(),
		"jti": tokenDB.ID,
		"nbf": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	// Simulate sending email
	fmt.Println("sent to:", user.Email)
	fmt.Printf("http://localhost:3000/login?token=%s\n", tokenString)

	// Return success
	render.Render(w, r, response.SuccessResponseOK(nil))
}

func (authApi *authApi) handlerValidateToken(w http.ResponseWriter, r *http.Request) {
	// Get token from query
	token := r.URL.Query().Get("token")
	if token == "" {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
		return
	}

	// Parse token
	claims := jwt.MapClaims{}
	_, errClaim := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

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
	tokenDB, err := authApi.DB.ClearUserToken(r.Context(), authStore.ClearUserTokenParams{
		TokenID: tokenUUID,
		Type:    "login",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	if errClaim != nil {
		render.Render(w, r, response.ErrorResponseUnauthorized(err))
		return
	}

	// Token valid
	// Create token to be sent for future requests
	tokenClaims := &authStore.SetUserTokenParams{
		UserID:     tokenDB.UserID,
		Type:       "access",
		Expiration: time.Now().Add(accessTokenExp),
	}
	tokenDB, err = authApi.DB.SetUserToken(r.Context(), *tokenClaims)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": tokenDB.ID,
		"exp": tokenDB.Expiration.Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"typ": tokenDB.Type,
	}).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	// Create token to be sent for future refresh
	tokenClaims = &authStore.SetUserTokenParams{
		UserID:     tokenDB.UserID,
		Type:       "refresh",
		Expiration: time.Now().Add(refreshTokenExp),
	}
	tokenDB, err = authApi.DB.SetUserToken(r.Context(), *tokenClaims)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": tokenDB.ID,
		"exp": tokenDB.Expiration.Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"typ": tokenDB.Type,
	}).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	render.Render(w, r, response.SuccessResponseOK(map[string]string{"accessToken": accessToken, "refreshToken": refreshToken}))
}

func (authApi *authApi) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get token from query
	token := r.URL.Query().Get("token")
	if token == "" {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidRefreshToken))
		return
	}

	// Parse token
	claims := jwt.MapClaims{}
	_, errClaims := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

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
	tokenDB, err := authApi.DB.ClearUserToken(r.Context(), authStore.ClearUserTokenParams{
		TokenID: tokenUUID,
		Type:    "refresh",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
			return
		}
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	_, err = authApi.DB.ClearUserToken(r.Context(), authStore.ClearUserTokenParams{
		UserID: tokenDB.UserID,
		Type:   "access",
	})
	if err != nil && err != sql.ErrNoRows {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	if errClaims != nil {
		render.Render(w, r, response.ErrorResponseUnauthorized(err))
		return
	}

	// Valid token
	// Create token to be sent for future requests
	tokenClaims := &authStore.SetUserTokenParams{
		UserID:     tokenDB.UserID,
		Type:       "access",
		Expiration: time.Now().Add(accessTokenExp),
	}
	tokenDB, err = authApi.DB.SetUserToken(r.Context(), *tokenClaims)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": tokenDB.ID,
		"exp": tokenDB.Expiration.Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"typ": tokenDB.Type,
	}).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	// Create token to be sent for future refresh
	tokenClaims = &authStore.SetUserTokenParams{
		UserID:     tokenDB.UserID,
		Type:       "refresh",
		Expiration: time.Now().Add(refreshTokenExp),
	}
	tokenDB, err = authApi.DB.SetUserToken(r.Context(), *tokenClaims)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": tokenDB.ID,
		"exp": tokenDB.Expiration.Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"typ": tokenDB.Type,
	}).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	render.Render(w, r, response.SuccessResponseOK(map[string]string{"accessToken": accessToken, "refreshToken": refreshToken}))
}

func (authApi *authApi) handlerLogout(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserKey).(authStore.User)
	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	_, err := authApi.DB.ClearAllTokenUser(r.Context(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}
	render.Render(w, r, response.SuccessResponseOK(nil))
}

func (authApi *authApi) handlerMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserKey).(authStore.User)

	if !ok {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	userResponse := struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		UpdateAt  time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}{
		ID:        user.ID,
		Name:      user.Name.String,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdatedAt,
		Email:     user.Email,
	}

	render.Render(w, r, response.SuccessResponseOK(userResponse))
}
