package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/auth"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	categoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	companyRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/health"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	inventoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	locationRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
	productRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed"
	supplierreturns "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns"
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
		auth.VersionCheck(redisClient),
	}

	r.Route("/api/v1", func(r chi.Router) {
		health.New(db).RegisterRoutes(r)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware...)

			product.New(db).RegisterRoutes(r)
			category.New(db).RegisterRoutes(r)
			company.New(db).RegisterRoutes(r)
			location.New(db).RegisterRoutes(r)
			inventory.New(db).RegisterRoutes(r)
			supplierreturns.New(db).RegisterRoutes(r)
		})

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
