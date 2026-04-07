package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := setupDatabase()

	defer db.Close()

	mux := http.NewServeMux()

	setupRoutes(mux, db)

	baseMiddleware := chain{recoverPanic, logger, secureHeaders, checkCors}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Fatal("SERVER_PORT environment variable is not set")
	}
	if err := http.ListenAndServe(":"+port, baseMiddleware.Then(mux)); err != nil {
		log.Fatal(err)
	}
}
