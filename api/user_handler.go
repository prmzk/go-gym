package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prmzk/go-base-prmzk/database"
	respond "github.com/prmzk/go-base-prmzk/json"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdatedAt,
	}
}

func (apiCfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsers(r.Context())
	if err != nil {
		respond.RespondWithError(w, 500, err.Error())
		return
	}

	if len(users) == 0 {
		respond.RespondWithError(w, 404, "No users found")
		return
	}

	var userList []User
	for _, u := range users {
		userList = append(userList, databaseUserToUser(u))
	}

	respond.RespondWithJSON(w, 200, userList)
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respond.RespondWithError(w, 400, "Invalid request payload")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		respond.RespondWithError(w, 500, "Internal server error")
		return
	}

	respond.RespondWithJSON(w, 200, databaseUserToUser(user))
}
