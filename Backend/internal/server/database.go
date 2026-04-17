package server

import (
	"log"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDatabase() *pgxpool.Pool {
	connString := os.Getenv("DB_HOST")
	if connString == "" {
		log.Fatal("DB environment variable is required")
	}

	db, err := database.NewPool(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
