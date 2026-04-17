package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	seedEnabled := flag.Bool("seed", false, "Enable seed endpoint")
	flag.Parse()

	if os.Getenv("DB_HOST") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db := server.SetupDatabase()

	defer db.Close()

	mux := http.NewServeMux()

	server.SetupRoutes(mux, db, *seedEnabled)

	baseMiddleware := server.Chain{server.RecoverPanic, server.Logger, server.SecureHeaders, server.CheckCORS}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Fatal("SERVER_PORT environment variable is not set")
	}
	if err := http.ListenAndServe(":"+port, baseMiddleware.Then(mux)); err != nil {
		log.Fatal(err)
	}
}
