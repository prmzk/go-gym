package api

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
	"github.com/prmzk/go-base-prmzk/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name.String,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdatedAt,
	}
}

func isValidEmail(email string) error {
	// Regular expression pattern for email validation
	// You can modify this pattern according to your requirements
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

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	body := &authRequest{}

	if err := render.Bind(r, body); err != nil {
		render.Render(w, r, ErrBadRequest(ErrInvalidEmail))
		return
	}

	// Check if email already exists
	_, err := apiCfg.DB.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			render.Render(w, r, ErrInternalServerError)
			return
		}
	} else {
		render.Render(w, r, ErrBadRequest(ErrDuplicateEmail))
		return
	}

	_, err = apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:    uuid.New(),
		Email: body.Email,
	})

	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, http.NoBody)
}

func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	body := &authRequest{}

	if err := render.Bind(r, body); err != nil {
		render.Render(w, r, ErrBadRequest(ErrInvalidEmail))
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, ErrBadRequest(ErrInvalidEmail))
			return
		}
		render.Render(w, r, ErrInternalServerError)
		return
	}

	tokenClaims := &database.SetUserTokenParams{
		ID:              user.ID,
		JwtID:           sql.NullString{String: uuid.New().String(), Valid: true},
		TokenExpiration: sql.NullTime{Time: time.Now().Add(time.Hour * 2), Valid: true},
	}

	updatedUser, err := apiCfg.DB.SetUserToken(r.Context(), *tokenClaims)

	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":               updatedUser.ID,
		"token_expiration": updatedUser.TokenExpiration.Time.Unix(),
		"jwt_id":           updatedUser.JwtID.String,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))

		return
	}

	fmt.Println("sent to:", updatedUser.Email)
	fmt.Printf("http://localhost:8080/v1/users/login/callback?token=%s\n", tokenString)

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, map[string]string{"message": "Email sent"})
}

func (apiCfg *apiConfig) handlerValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
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
		render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
		return
	}
	fmt.Println(claims)

	userUUID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
		return
	}

	user, err := apiCfg.DB.GetUserById(r.Context(), userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, ErrInternalServerError)
			return
		}
		render.Render(w, r, ErrUnauthorized(ErrInvalidBearerToken))
		return
	}

	if user.JwtID.String != claims["jwt_id"].(string) {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	if user.TokenExpiration.Time.Unix() != int64(claims["token_expiration"].(float64)) {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	if user.TokenExpiration.Time.Before(time.Now().UTC()) {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	clearedUser, err := apiCfg.DB.ClearUserToken(r.Context(), user.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	authToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": clearedUser.ID,
	}).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, struct {
		Token string `json:"token"`
	}{
		Token: authToken,
	})
}

func (apiCfg *apiConfig) handlerMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userKey).(database.User)

	if !ok {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, databaseUserToUser(user))
}
