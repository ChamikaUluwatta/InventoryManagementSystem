package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	seedEnabled := flag.Bool("seed", false, "Enable seed endpoint")
	flag.Parse()

	db := setupDatabase()

	mux := http.NewServeMux()

	setupRoutes(mux, db, seedEnabled)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	fmt.Printf("Server starting on :%s\n", os.Getenv("DB_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("DB_PORT"), checkCors(mux)))
}

func setupRoutes(mux *http.ServeMux, db *pgxpool.Pool, seedEnabled *bool) {
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

	if *seedEnabled {
		fmt.Println("Seed endpoint is registered.")
		seedService := seed.NewService(companyRepo, categoryRepo, locationRepo, productRepo, inventoryRepo, db)
		seedHandler := seed.NewHandler(seedService)
		seedHandler.RegisterRoutes(mux)
	}

	productHandler.RegisterRoutes(mux)

	categoryHandler.RegisterRoutes(mux)

	companyHandler.RegisterRoutes(mux)

	locationHandler.RegisterRoutes(mux)

	inventoryHandler.RegisterRoutes(mux)
}

func setupDatabase() *pgxpool.Pool {
	connString := os.Getenv("DB_HOST")
	if connString == "" {
		log.Fatal("DB environment variable is required")
	}

	db, err := database.NewPool(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	return db
}

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}

var allowedList = []string{
	"http://localhost:5173",
	"http://127.0.0.1:5173",
}

var allowedMethods = []string{
	"GET",
	"DELETE",
	"PUT",
	"POST",
	"OPTIONS",
	"UPDATE",
}

func checkCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPreflight(r) {
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			if slices.Contains(allowedList, origin) && slices.Contains(allowedMethods, method) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
			}
		} else {
			origin := r.Header.Get("Origin")
			if slices.Contains(allowedList, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		w.Header().Add("Vary", "Origin")
		next.ServeHTTP(w, r)
	})
}
