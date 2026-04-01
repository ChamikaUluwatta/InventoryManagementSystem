package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, inventory *Inventory) error
	GetByID(ctx context.Context, id int) (*Inventory, error)
	GetAll(ctx context.Context) ([]Inventory, error)
	Update(ctx context.Context, inventory *Inventory) error
	Delete(ctx context.Context, id int) error
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error)
	GetByLocation(ctx context.Context, locationID string) ([]Inventory, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, inventory *Inventory) error {
	query := `
		INSERT INTO "inventories" (product_id, location_id, stock)
		VALUES (@product_id, @location_id, @stock)
		RETURNING inventory_id`

	args := pgx.NamedArgs{
		"product_id":  inventory.ProductID,
		"location_id": inventory.LocationID,
		"stock":       inventory.Stock,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(&inventory.InventoryID)

	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*Inventory, error) {
	query := `
		SELECT inventory_id, product_id, location_id, stock
		FROM "inventories"
		WHERE inventory_id = @inventory_id`

	var inventory Inventory
	args := pgx.NamedArgs{
		"inventory_id": id,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(
		&inventory.InventoryID,
		&inventory.ProductID,
		&inventory.LocationID,
		&inventory.Stock,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by id: %w", err)
	}
	return &inventory, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Inventory, error) {
	query := `
		SELECT inventory_id, product_id, location_id, stock
		FROM "inventories"
		ORDER BY inventory_id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all inventories: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Inventory])
}

func (r *repository) Update(ctx context.Context, inventory *Inventory) error {
	query := `
		UPDATE "inventories"
		SET product_id = @product_id, location_id = @location_id, stock = @stock
		WHERE inventory_id = @inventory_id`

	args := pgx.NamedArgs{
		"product_id":   inventory.ProductID,
		"location_id":  inventory.LocationID,
		"stock":        inventory.Stock,
		"inventory_id": inventory.InventoryID,
	}
	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "inventories" WHERE inventory_id = @inventory_id`
	args := pgx.NamedArgs{
		"inventory_id": id,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}
	return nil
}

func (r *repository) GetByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error) {
	query := `
		SELECT inventory_id, product_id, location_id, stock
		FROM "inventories"
		WHERE product_id = @product_id
		ORDER BY inventory_id`

	args := pgx.NamedArgs{
		"product_id": productID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventories by product: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Inventory])
}

func (r *repository) GetByLocation(ctx context.Context, locationID string) ([]Inventory, error) {
	query := `
		SELECT inventory_id, product_id, location_id, stock
		FROM "inventories"
		WHERE location_id = @location_id
		ORDER BY inventory_id`

	args := pgx.NamedArgs{
		"location_id": locationID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventories by location: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Inventory])
}
