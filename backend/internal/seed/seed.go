package seed

import (
	categoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	companyRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	inventoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	locationRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	productRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *handler.Handler {
	svc := service.NewService(
		companyRepo.NewRepository(db),
		categoryRepo.NewRepository(db),
		locationRepo.NewRepository(db),
		productRepo.NewRepository(db),
		inventoryRepo.NewRepository(db),
		db,
	)
	return handler.NewHandler(svc)
}
