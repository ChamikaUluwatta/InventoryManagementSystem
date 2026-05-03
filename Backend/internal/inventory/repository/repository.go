package repository

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, inventory *model.Inventory) error
	GetByID(ctx context.Context, id int) (*model.Inventory, error)
	GetAll(ctx context.Context) ([]model.Inventory, error)
	Update(ctx context.Context, inventory *model.Inventory) error
	Delete(ctx context.Context, id int) error
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]model.Inventory, error)
	GetByLocation(ctx context.Context, locationID string) ([]model.Inventory, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, inventory *model.Inventory) error {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return apperror.BadRequest("invalid product_id or location_id", err)
			case "23505":
				return apperror.Conflict("inventory already exists for product and location", err)
			}
		}
		return apperror.Internal("failed to create inventory", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.Inventory, error) {
	query := `
		SELECT inventory_id, product_id, location_id, stock
		FROM "inventories"
		WHERE inventory_id = @inventory_id`

	var inventory model.Inventory
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("inventory not found", err)
		}
		return nil, apperror.Internal("failed to get inventory by id", err)
	}
	return &inventory, nil
}

func (r *repository) GetAll(ctx context.Context) ([]model.Inventory, error) {
	query := `
		SELECT inventory_id, product_id, location_id, stock
		FROM "inventories"
		ORDER BY inventory_id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, apperror.Internal("failed to get all inventories", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Inventory])
}

func (r *repository) Update(ctx context.Context, inventory *model.Inventory) error {
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
	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return apperror.BadRequest("invalid product_id or location_id", err)
			case "23505":
				return apperror.Conflict("inventory already exists for product and location", err)
			}
		}
		return apperror.Internal("failed to update inventory", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("inventory not found", nil)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "inventories" WHERE inventory_id = @inventory_id`
	args := pgx.NamedArgs{
		"inventory_id": id,
	}

	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to delete inventory", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("inventory not found", nil)
	}
	return nil
}

func (r *repository) GetByProduct(ctx context.Context, productID uuid.UUID) ([]model.Inventory, error) {
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
		return nil, apperror.Internal("failed to get inventories by product", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Inventory])
}

func (r *repository) GetByLocation(ctx context.Context, locationID string) ([]model.Inventory, error) {
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
		return nil, apperror.Internal("failed to get inventories by location", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Inventory])
}