package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prmzk/go-base-prmzk/api"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		godotenv.Load()
	}

	r, err := api.NewRouter(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Go Gym Server is running on port 8080...")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
