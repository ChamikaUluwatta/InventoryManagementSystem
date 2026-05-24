package health

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/health/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/health/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *handler.Handler {
	dbChecker := service.NewDatabaseHealthChecker(db)
	svc := service.NewService(dbChecker)
	return handler.NewHandler(svc)
}
