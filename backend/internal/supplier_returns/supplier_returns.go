package supplierreturns

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *handler.Handler {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	return handler.NewHandler(svc)
}
