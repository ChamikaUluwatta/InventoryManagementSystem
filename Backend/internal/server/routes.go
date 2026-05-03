package server

import (
	"fmt"
	"net/http"

	categoryHandler "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/handler"
	categoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	categorySvc "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/service"
	companyHandler "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/handler"
	companyRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	companySvc "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/service"
	inventoryHandler "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/handler"
	inventoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	inventorySvc "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/service"
	locationHandler "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/handler"
	locationRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	locationSvc "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/service"
	productHandler "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/handler"
	productRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	productSvc "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed"
	supplierReturnsHandler "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/handler"
	supplierReturnsRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/repository"
	supplierReturnsSvc "github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(mux *http.ServeMux, db *pgxpool.Pool, seedEnabled bool) {
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	productRepoInstance := productRepo.NewRepository(db)
	productService := productSvc.NewService(productRepoInstance)
	productHandlerInstance := productHandler.NewHandler(productService)

	categoryRepoInstance := categoryRepo.NewRepository(db)
	categoryService := categorySvc.NewService(categoryRepoInstance)
	categoryHandlerInstance := categoryHandler.NewHandler(categoryService)

	companyRepoInstance := companyRepo.NewRepository(db)
	companyService := companySvc.NewService(companyRepoInstance)
	companyHandlerInstance := companyHandler.NewHandler(companyService)

	locationRepoInstance := locationRepo.NewRepository(db)
	locationService := locationSvc.NewService(locationRepoInstance)
	locationHandlerInstance := locationHandler.NewHandler(locationService)

	inventoryRepoInstance := inventoryRepo.NewRepository(db)
	inventoryService := inventorySvc.NewService(inventoryRepoInstance)
	inventoryHandlerInstance := inventoryHandler.NewHandler(inventoryService)

	supplierReturnsRepoInstance := supplierReturnsRepo.NewRepository(db)
	supplierReturnsService := supplierReturnsSvc.NewService(supplierReturnsRepoInstance)
	supplierReturnsHandlerInstance := supplierReturnsHandler.NewHandler(supplierReturnsService)

	if seedEnabled {
		fmt.Println("Seed endpoint is registered.")
		seedService := seed.NewService(companyRepoInstance, categoryRepoInstance, locationRepoInstance, productRepoInstance, inventoryRepoInstance, db)
		seedHandler := seed.NewHandler(seedService)
		seedHandler.RegisterRoutes(mux)
	}

	productHandlerInstance.RegisterRoutes(mux)
	categoryHandlerInstance.RegisterRoutes(mux)
	companyHandlerInstance.RegisterRoutes(mux)
	locationHandlerInstance.RegisterRoutes(mux)
	inventoryHandlerInstance.RegisterRoutes(mux)
	supplierReturnsHandlerInstance.RegisterRoutes(mux)
}