package inventory

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *handler.Handler {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	return handler.NewHandler(svc)
}
