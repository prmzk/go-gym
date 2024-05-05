package auth

import (
	"database/sql"
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

func isValidEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	if !match {
		return ErrInvalidEmail
	}
	return nil
}

type authRequest struct {
	Email string
}

func (body *authRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	return isValidEmail(body.Email)
}

func (authApi *authApi) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	body := &authRequest{}
	if err := render.Bind(r, body); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(ErrInvalidEmail))
		return
	}

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

	_, err = authApi.DB.CreateUser(r.Context(), authStore.CreateUserParams{
		ID:    uuid.New(),
		Email: body.Email,
	})

	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	render.Render(w, r, response.SuccessResponseCreated(nil))
}

func (authApi *authApi) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	body := &authRequest{}

	if err := render.Bind(r, body); err != nil {
		render.Render(w, r, response.ErrorResponseBadRequest(ErrInvalidEmail))
		return
	}

	user, err := authApi.DB.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrorResponseBadRequest(ErrInvalidEmail))
			return
		}
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	tokenClaims := &authStore.SetUserTokenParams{
		ID:              user.ID,
		JwtID:           sql.NullString{String: uuid.New().String(), Valid: true},
		TokenExpiration: sql.NullTime{Time: time.Now().Add(time.Hour * 2), Valid: true},
	}

	updatedUser, err := authApi.DB.SetUserToken(r.Context(), *tokenClaims)

	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":               updatedUser.ID,
		"token_expiration": updatedUser.TokenExpiration.Time.Unix(),
		"jwt_id":           updatedUser.JwtID.String,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))

		return
	}

	fmt.Println("sent to:", updatedUser.Email)
	fmt.Printf("http://localhost:8080/v1/users/login/callback?token=%s\n", tokenString)

	render.Render(w, r, response.SuccessResponseOK(nil))
}

func (authApi *authApi) handlerValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
		return
	}

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

	userUUID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
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

	if user.JwtID.String != claims["jwt_id"].(string) {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
		return
	}

	if user.TokenExpiration.Time.Unix() != int64(claims["token_expiration"].(float64)) {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
		return
	}

	if user.TokenExpiration.Time.Before(time.Now().UTC()) {
		render.Render(w, r, response.ErrorResponseUnauthorized(ErrInvalidBearerToken))
		return
	}

	clearedUser, err := authApi.DB.ClearUserToken(r.Context(), user.ID)
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	authToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JwtSub: clearedUser.ID,
	}).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, response.ErrorResponseInternalServerError())
		return
	}

	render.Render(w, r, response.SuccessResponseOK(map[string]string{"token": authToken}))
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
	}{
		ID:        user.ID,
		Name:      user.Name.String,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdatedAt,
	}

	render.Render(w, r, response.SuccessResponseOK(userResponse))
}
