package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/database"
	respond "github.com/prmzk/go-base-prmzk/json"
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

func isValidEmail(email string) bool {
	// Regular expression pattern for email validation
	// You can modify this pattern according to your requirements
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respond.RespondWithError(w, 400, "Invalid request payload")
		return
	}

	if params.Email == "" {
		respond.RespondWithError(w, 400, "Email cannot be empty")
		return
	}

	// Tidy up the email string
	params.Email = strings.TrimSpace(params.Email)
	params.Email = strings.ToLower(params.Email)

	// Validate email format
	if !isValidEmail(params.Email) {
		respond.RespondWithError(w, 400, "Invalid email format")
		return
	}

	// Check if email already exists
	_, err = apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			respond.RespondWithError(w, 500, err.Error())
			return
		}
	} else {
		respond.RespondWithError(w, 400, fmt.Sprintf("Email %s already exists", params.Email))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:    uuid.New(),
		Email: params.Email,
	})

	if err != nil {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	respond.RespondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respond.RespondWithError(w, 400, "Invalid request payload")
		return
	}

	if params.Email == "" {
		respond.RespondWithError(w, 400, "Email cannot be empty")
		return
	}

	// Tidy up the email string
	params.Email = strings.TrimSpace(params.Email)
	params.Email = strings.ToLower(params.Email)

	// Validate email format
	if !isValidEmail(params.Email) {
		respond.RespondWithError(w, 400, "Invalid email format")
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.RespondWithError(w, 404, "User not found")
			return
		}
		respond.RespondWithError(w, 500, err.Error())
		return
	}

	tokenClaims := &database.SetUserTokenParams{
		ID:              user.ID,
		JwtID:           sql.NullString{String: uuid.New().String(), Valid: true},
		TokenExpiration: sql.NullTime{Time: time.Now().Add(time.Hour * 2), Valid: true},
	}

	updatedUser, err := apiCfg.DB.SetUserToken(r.Context(), *tokenClaims)

	if err != nil {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":               updatedUser.ID,
		"token_expiration": updatedUser.TokenExpiration.Time.Unix(),
		"jwt_id":           updatedUser.JwtID.String,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	fmt.Println("sent to:", updatedUser.Email)
	fmt.Printf("http://localhost:8080/v1/users/login/callback?token=%s\n", tokenString)
	respond.RespondWithJSON(w, 200, map[string]string{"message": "Email sent"})
}

func (apiCfg *apiConfig) handlerValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respond.RespondWithError(w, 400, "Token not found in query string")
		return
	}

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
		respond.RespondWithError(w, 401, "Invalid token")
		return
	}
	fmt.Println(claims)

	// Get the user ID from the token claims
	userID, ok := claims["id"].(string)
	if !ok {
		respond.RespondWithError(w, 401, "Invalid token")
		return
	}

	expiration, ok := claims["token_expiration"].(float64)
	if !ok {
		respond.RespondWithError(w, 401, "Invalid token")
		return
	}

	jwtID, ok := claims["jwt_id"].(string)
	if !ok {
		respond.RespondWithError(w, 401, "Invalid token")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respond.RespondWithError(w, 400, "Invalid user ID")
		return
	}

	user, err := apiCfg.DB.GetUserById(r.Context(), userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.RespondWithError(w, 404, "User not found")
			return
		}
		respond.RespondWithError(w, 500, err.Error())
		return
	}

	if user.JwtID.String != jwtID {
		respond.RespondWithError(w, 401, "Invalid token")
		return
	}

	if user.TokenExpiration.Time.Unix() != int64(expiration) {
		respond.RespondWithError(w, 401, "Invalid token")
		return
	}

	if user.TokenExpiration.Time.Before(time.Now().UTC()) {
		respond.RespondWithError(w, 401, "Token expired. Please login again.")
		return
	}

	clearedUser, err := apiCfg.DB.ClearUserToken(r.Context(), user.ID)
	if err != nil {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": clearedUser.ID,
	})
	authTokenString, err := authToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	respond.RespondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{
		Token: authTokenString,
	})
}

func (apiCfg *apiConfig) handlerMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userKey).(database.User)
	fmt.Println(user)

	if !ok {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	respond.RespondWithJSON(w, 200, databaseUserToUser(user))
}
