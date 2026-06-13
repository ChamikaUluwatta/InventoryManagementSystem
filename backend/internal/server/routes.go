package server

import (
	"net/http"
	"os"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/auth"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/health"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
	// "github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed"
	supplierreturns "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r chi.Router, db *pgxpool.Pool) {

	jwksURL := os.Getenv("AUTH_JWKS_URL")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	jwks := auth.NewJWKSCache(jwksURL)
	redisClient := auth.NewRedisClient(redisHost, redisPort, redisPassword)
	authMiddleware := []func(http.Handler) http.Handler{
		auth.JWTAuth(jwks),
		auth.VersionCheck(redisClient),
		TenantRLS(db),
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
			// seed.New(db).RegisterRoutes(r)
		})
	})
}
