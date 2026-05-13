package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthChecker interface {
	Name() string
	Check(ctx context.Context) error
}

type CheckResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type DatabaseHealthChecker struct {
	db *pgxpool.Pool
}

func NewDatabaseHealthChecker(db *pgxpool.Pool) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

func (d *DatabaseHealthChecker) Name() string {
	return "Database"
}

func (d *DatabaseHealthChecker) Check(ctx context.Context) error {
	if err := d.db.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}
