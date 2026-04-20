package server

import (
	"fmt"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed"
	supplierreturns "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(mux *http.ServeMux, db *pgxpool.Pool, seedEnabled bool) {
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

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

	supplierReturnsRepo := supplierreturns.NewRepository(db)
	supplierReturnsService := supplierreturns.NewService(supplierReturnsRepo)
	supplierReturnsHandler := supplierreturns.NewHandler(supplierReturnsService)

	if seedEnabled {
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
	supplierReturnsHandler.RegisterRoutes(mux)
}
