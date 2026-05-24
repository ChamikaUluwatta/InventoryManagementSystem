package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/auth"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/health"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed"
	supplierreturns "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns"
	categoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	companyRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	inventoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	locationRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	productRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r chi.Router, db *pgxpool.Pool, seedEnabled bool) {

	// Initialize auth middleware
	jwksURL := os.Getenv("AUTH_JWKS_URL")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	jwks := auth.NewJWKSCache(jwksURL)
	redisClient := auth.NewRedisClient(redisHost, redisPort)
	authMiddleware := []func(http.Handler) http.Handler{
		auth.JWTAuth(jwks),
		auth.PermissionCheck(redisClient),
	}

	r.Route("/api/v1", func(r chi.Router) {
		health.New(db).RegisterRoutes(r)
		product.New(db).RegisterRoutes(r, authMiddleware...)
		category.New(db).RegisterRoutes(r, authMiddleware...)
		company.New(db).RegisterRoutes(r, authMiddleware...)
		location.New(db).RegisterRoutes(r, authMiddleware...)
		inventory.New(db).RegisterRoutes(r, authMiddleware...)
		supplierreturns.New(db).RegisterRoutes(r, authMiddleware...)

		if seedEnabled {
			fmt.Println("Seed endpoint is registered.")
			seedService := seed.NewService(
				companyRepo.NewRepository(db),
				categoryRepo.NewRepository(db),
				locationRepo.NewRepository(db),
				productRepo.NewRepository(db),
				inventoryRepo.NewRepository(db),
				db,
			)
			seedHandler := seed.NewHandler(seedService)
			seedHandler.RegisterRoutes(r)
		}
	})
}
