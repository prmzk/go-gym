package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/prmzk/go-base-prmzk/api/response"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, response.SuccessResponseOK(nil))
}
