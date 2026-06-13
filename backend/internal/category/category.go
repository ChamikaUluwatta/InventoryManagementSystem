package category

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *handler.Handler {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	return handler.NewHandler(svc)
}
