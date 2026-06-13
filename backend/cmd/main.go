package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("DB_HOST") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db := server.SetupDatabase()

	defer db.Close()

	r := chi.NewRouter()

	r.Use(server.RecoverPanic)
	r.Use(server.Logger)
	r.Use(server.SecureHeaders)
	// r.Use(server.CheckCORS)

	server.SetupRoutes(r, db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Fatal("SERVER_PORT environment variable is not set")
	}
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
