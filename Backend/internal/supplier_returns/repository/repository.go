package repository

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierReturnRepository interface {
	Create(ctx context.Context, req *model.CreateSupplierReturnRequest) (*model.SupplierReturn, error)
	GetByID(ctx context.Context, id int) (*model.SupplierReturn, error)
	GetAll(ctx context.Context) ([]model.SupplierReturn, error)
	UpdateStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) SupplierReturnRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, req *model.CreateSupplierReturnRequest) (*model.SupplierReturn, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperror.Internal("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	headerQuery := `
		INSERT INTO "supplier_returns" (return_no, company_id, status, reason, notes)
		VALUES (@return_no, @company_id, @status, @reason, @notes)
		RETURNING supplier_return_id, created_at, approved_at, completed_at`

	var created model.SupplierReturn

	err = tx.QueryRow(ctx, headerQuery, pgx.NamedArgs{
		"return_no":  req.ReturnNo,
		"company_id": req.CompanyID,
		"status":     model.ReturnStatusDraft,
		"reason":     req.Reason,
		"notes":      req.Notes,
	}).Scan(
		&created.SupplierReturnID,
	)
	if err != nil {
		return nil, toAppError("failed to create supplier return", err)
	}

	itemQuery := `
		INSERT INTO "supplier_return_items"
			(supplier_return_id, product_id, location_id, quantity, unit_cost, product_name_snapshot, location_snapshot)
		VALUES
			(@supplier_return_id, @product_id, @location_id, @quantity, @unit_cost, @product_name_snapshot, @location_snapshot)
		RETURNING supplier_return_item_id, supplier_return_id, product_id, location_id, quantity, unit_cost, product_name_snapshot, location_snapshot`

	snapshotQuery := `
		SELECT p.product_name, l.location_id
		FROM "products" p
		JOIN "locations" l ON l.location_id = @location_id
		WHERE p.product_id = @product_id
		  AND p.company_id = @company_id`

	for _, item := range req.Items {
		var productNameSnapshot string
		var locationSnapshot string
		err := tx.QueryRow(ctx, snapshotQuery, pgx.NamedArgs{
			"product_id":  *item.ProductID,
			"location_id": *item.LocationID,
			"company_id":  req.CompanyID,
		}).Scan(&productNameSnapshot, &locationSnapshot)
		if err != nil {
			return nil, apperror.Internal("Internal server error", err)
		}

		_, err = tx.Query(ctx, itemQuery, pgx.NamedArgs{
			"supplier_return_id":    created.SupplierReturnID,
			"product_id":            *item.ProductID,
			"location_id":           *item.LocationID,
			"quantity":              item.Quantity,
			"unit_cost":             item.UnitCost,
			"product_name_snapshot": productNameSnapshot,
			"location_snapshot":     locationSnapshot,
		})
		if err != nil {
			return nil, toAppError("Internal Server Error", err)
		}

	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperror.Internal("Internal server error", err)
	}

	return &created, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.SupplierReturn, error) {
	headerQuery := `
		SELECT supplier_return_id, return_no, company_id, status, reason, notes, created_at, approved_at, completed_at
		FROM "supplier_returns"
		WHERE supplier_return_id = @supplier_return_id`

	headRows, err := r.db.Query(ctx, headerQuery, pgx.NamedArgs{
		"supplier_return_id": id,
	})
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	result, err := pgx.CollectOneRow(headRows, pgx.RowToStructByName[model.SupplierReturn])
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	itemsQuery := `
		SELECT supplier_return_item_id, supplier_return_id, product_id, location_id, quantity, unit_cost, product_name_snapshot, location_snapshot
		FROM "supplier_return_items"
		WHERE supplier_return_id = @supplier_return_id
		ORDER BY supplier_return_item_id`

	itemRows, err := r.db.Query(ctx, itemsQuery, pgx.NamedArgs{
		"supplier_return_id": id,
	})
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	items, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[model.SupplierReturnItem])
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	result.Items = items
	return &result, nil
}

func (r *repository) GetAll(ctx context.Context) ([]model.SupplierReturn, error) {
	query := `
		SELECT supplier_return_id, return_no, company_id, status, created_at, approved_at, completed_at
		FROM "supplier_returns"
		ORDER BY created_at DESC, supplier_return_id DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}
	defer rows.Close()

	returns, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.SupplierReturn])
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return returns, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
	query := `
		UPDATE "supplier_returns"
		SET
			status = @status,
			approved_at = CASE
				WHEN @status = 'approved' AND approved_at IS NULL THEN NOW()
				ELSE approved_at
			END,
			completed_at = CASE
				WHEN @status IN ('completed', 'credited') AND completed_at IS NULL THEN NOW()
				ELSE completed_at
			END
		WHERE supplier_return_id = @supplier_return_id
		RETURNING supplier_return_id, return_no, company_id, status, reason, notes, created_at, approved_at, completed_at`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{
		"supplier_return_id": id,
		"status":             status,
	})
	if err != nil {
		return nil, toAppError("failed to update supplier return status", err)
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.SupplierReturn])
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return &updated, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "supplier_returns" WHERE supplier_return_id = @supplier_return_id`

	result, err := r.db.Exec(ctx, query, pgx.NamedArgs{
		"supplier_return_id": id,
	})
	if err != nil {
		return toAppError("failed to delete supplier return", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("supplier return not found", nil)
	}

	return nil
}

func toAppError(message string, err error) error {
	return apperror.Internal(message, err)
}