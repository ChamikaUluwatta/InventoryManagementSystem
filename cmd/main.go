package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := database.NewPool(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	companyRepo := company.NewRepository(db)
	companyService := company.NewService(companyRepo)
	companyHandler := company.NewHandler(companyService)

	locationRepo := location.NewRepository(db)
	locationService := location.NewService(locationRepo)
	locationHandler := location.NewHandler(locationService)

	inventoryRepo := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepo)
	inventoryHandler := inventory.NewHandler(inventoryService)

	mux := http.NewServeMux()

	productHandler.RegisterRoutes(mux)

	categoryHandler.RegisterRoutes(mux)

	companyHandler.RegisterRoutes(mux)

	locationHandler.RegisterRoutes(mux)

	inventoryHandler.RegisterRoutes(mux)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
