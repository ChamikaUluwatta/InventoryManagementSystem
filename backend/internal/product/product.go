package product

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *handler.Handler {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	return handler.NewHandler(svc)
}
