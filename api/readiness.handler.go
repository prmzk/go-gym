package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, http.NoBody)
}
