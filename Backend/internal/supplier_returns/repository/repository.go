package repository

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierReturnRepository interface {
	Create(ctx context.Context, req *model.SupplierReturn) error
	GetByID(ctx context.Context, id int) (*model.SupplierReturn, error)
	GetAll(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error)
	UpdateStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) SupplierReturnRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, req *model.SupplierReturn) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.Internal("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	headerQuery := `
		INSERT INTO "supplier_returns" (return_no, company_id, status, reason, notes)
		VALUES (@return_no, @company_id, @status, @reason, @notes)
		RETURNING supplier_return_id, created_at, approved_at, completed_at`

	err = tx.QueryRow(ctx, headerQuery, pgx.NamedArgs{
		"return_no":  req.ReturnNo,
		"company_id": req.CompanyID,
		"status":     model.ReturnStatusDraft,
		"reason":     req.Reason,
		"notes":      req.Notes,
	}).Scan(
		&req.SupplierReturnID,
		&req.CreatedAt,
		&req.ApprovedAt,
		&req.CompletedAt,
	)
	if err != nil {
		return apperror.Internal("failed to create supplier return", err)
	}

	itemQuery := `
    INSERT INTO "supplier_return_items"
        (supplier_return_id, product_id, location_id, quantity, unit_cost, product_name_snapshot, location_snapshot)
    SELECT @supplier_return_id, @product_id, @location_id, @quantity, @unit_cost,
           p.product_name, @location_id
    FROM "products" p
    WHERE p.product_id = @product_id AND p.company_id = @company_id
    RETURNING supplier_return_item_id,product_name_snapshot, location_snapshot`

	batch := pgx.Batch{}
	for _, item := range req.Items {
		batch.Queue(itemQuery, pgx.NamedArgs{
			"supplier_return_id": req.SupplierReturnID,
			"product_id":         item.ProductID,
			"location_id":        item.LocationID,
			"quantity":           item.Quantity,
			"unit_cost":          item.UnitCost,
			"company_id":         req.CompanyID,
		})
	}

	br := tx.SendBatch(ctx, &batch)

	items := req.Items
	req.Items = make([]model.SupplierReturnItem, 0, len(items))
	for _, item := range items {
		var insertedItem = model.SupplierReturnItem{
			SupplierReturnID: item.SupplierReturnID,
			ProductID:        item.ProductID,
			LocationID:       item.LocationID,
			Quantity:         item.Quantity,
			UnitCost:         item.UnitCost,
		}
		if err := br.QueryRow().Scan(
			&insertedItem.SupplierReturnItemID,
			&insertedItem.ProductNameSnapshot,
			&insertedItem.LocationSnapshot,
		); err != nil {
			br.Close()
			if errors.Is(err, pgx.ErrNoRows) {
				return apperror.BadRequest("invalid product_id or location_id for item", nil)
			}
			return apperror.Internal("failed to insert supplier return item", err)
		}
		req.Items = append(req.Items, insertedItem)
	}

	br.Close()

	if err := tx.Commit(ctx); err != nil {
		return apperror.Internal("Internal server error", err)
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.SupplierReturn, error) {
	query := `
		SELECT
			sr.supplier_return_id, sr.return_no, sr.company_id, sr.status, sr.reason, sr.notes, sr.created_at, sr.approved_at, sr.completed_at,
			sri.supplier_return_item_id, sri.product_id, sri.location_id, sri.quantity, sri.unit_cost, sri.product_name_snapshot, sri.location_snapshot
		FROM "supplier_returns" sr
		LEFT JOIN "supplier_return_items" sri ON sri.supplier_return_id = sr.supplier_return_id
		WHERE sr.supplier_return_id = @supplier_return_id
		ORDER BY sri.supplier_return_item_id`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{
		"supplier_return_id": id,
	})
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}
	defer rows.Close()

	var result *model.SupplierReturn
	for rows.Next() {
		var header model.SupplierReturn
		var item model.SupplierReturnItem
		if err := rows.Scan(
			&header.SupplierReturnID, &header.ReturnNo, &header.CompanyID, &header.Status, &header.Reason, &header.Notes,
			&header.CreatedAt, &header.ApprovedAt, &header.CompletedAt,
			&item.SupplierReturnItemID, &item.ProductID, &item.LocationID, &item.Quantity, &item.UnitCost,
			&item.ProductNameSnapshot, &item.LocationSnapshot,
		); err != nil {
			return nil, apperror.Internal("Internal Server Error", err)
		}

		if result == nil {
			header.Items = []model.SupplierReturnItem{item}
			result = &header
		} else {
			result.Items = append(result.Items, item)
		}
	}

	if result == nil {
		return nil, apperror.NotFound("supplier return not found", nil)
	}

	return result, nil
}

func (r *repository) GetAll(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
	query := `
		SELECT 
			supplier_return_id,
			return_no, company_id, 
			status, 
			reason, 
			notes, 
			created_at, 
			approved_at, 
			completed_at
		FROM "supplier_returns"
		ORDER BY created_at DESC, supplier_return_id DESC
		LIMIT @limit OFFSET @offset`

	args := pgx.NamedArgs{
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	rows, err := r.db.Query(ctx, query, args)
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
				WHEN @status::VARCHAR = 'approved' AND approved_at IS NULL THEN NOW()
				ELSE approved_at
			END,
			completed_at = CASE
				WHEN @status::VARCHAR IN ('completed', 'credited') AND completed_at IS NULL THEN NOW()
				ELSE completed_at
			END
		WHERE supplier_return_id = @supplier_return_id
		RETURNING supplier_return_id, return_no, company_id, status, reason, notes, created_at, approved_at, completed_at`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{
		"supplier_return_id": id,
		"status":             string(status),
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
