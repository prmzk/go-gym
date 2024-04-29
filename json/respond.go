package json

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Internal Server Error: ", message)

	}

	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errResponse{Error: message})
}
