package testutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestDB struct {
	Pool      *pgxpool.Pool
	container *postgres.PostgresContainer
}

func SetupTestDB(ctx context.Context, migrationsDir string) (*TestDB, error) {
	container, err := StartPostgresContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	pool, err := NewPoolFromContainer(ctx, container)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("new pool: %w", err)
	}

	if err := RunMigrations(ctx, pool, migrationsDir); err != nil {
		pool.Close()
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &TestDB{
		Pool:      pool,
		container: container,
	}, nil
}

func (tdb *TestDB) Close() {
	tdb.Pool.Close()
	_ = tdb.container.Terminate(context.Background())
}
