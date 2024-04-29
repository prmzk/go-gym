package api

import (
	"net/http"

	"github.com/prmzk/go-base-prmzk/json"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	json.RespondWithJSON(w, 200, struct{}{})
}
