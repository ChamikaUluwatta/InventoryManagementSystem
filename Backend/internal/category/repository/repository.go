package repository

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int) (*model.Category, error)
	GetAll(ctx context.Context, params model.QueryParams) ([]model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int) error
	GetByParent(ctx context.Context, parentID *int) ([]model.Category, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) db(ctx context.Context) database.Querier {
	if tx, ok := database.GetTx(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *repository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO "categories" (category_name, parent_id, tenant_id)
		VALUES (@category_name, @parent_id, @tenant_id)
		RETURNING category_id`

	args := pgx.NamedArgs{
		"category_name": category.CategoryName,
		"parent_id":     category.ParentID,
		"tenant_id":     category.TenantID,
	}
	err := r.db(ctx).QueryRow(ctx, query, args).Scan(&category.CategoryID)

	if err != nil {
		return apperror.Internal("failed to create category", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	query := `
		SELECT category_id, category_name, parent_id
		FROM "categories"
		WHERE category_id = @category_id`

	var category model.Category
	args := pgx.NamedArgs{
		"category_id": id,
	}
	err := r.db(ctx).QueryRow(ctx, query, args).Scan(
		&category.CategoryID,
		&category.CategoryName,
		&category.ParentID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("category not found", err)
		}
		return nil, apperror.Internal("failed to get category by id", err)
	}
	return &category, nil
}

func (r *repository) GetAll(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
	query := `
		SELECT category_id, category_name, parent_id
		FROM "categories"
		ORDER BY category_name
		LIMIT @limit OFFSET @offset`

	args := pgx.NamedArgs{
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	rows, err := r.db(ctx).Query(ctx, query, args)
	if err != nil {
		return nil, apperror.Internal("failed to get all categories", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
}

func (r *repository) Update(ctx context.Context, category *model.Category) error {
	query := `
		UPDATE "categories"
		SET category_name = @category_name, parent_id = @parent_id
		WHERE category_id = @category_id`
	args := pgx.NamedArgs{
		"category_name": category.CategoryName,
		"parent_id":     category.ParentID,
		"category_id":   category.CategoryID,
	}
	result, err := r.db(ctx).Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to update category", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("category not found", nil)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "categories" WHERE category_id = @category_id`
	args := pgx.NamedArgs{
		"category_id": id,
	}

	result, err := r.db(ctx).Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to delete category", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("category not found", nil)
	}
	return nil
}

func (r *repository) GetByParent(ctx context.Context, parentID *int) ([]model.Category, error) {
	query := `
		SELECT category_id, category_name, parent_id
		FROM "categories"
		WHERE parent_id = @parent_id
		ORDER BY category_name`

	args := pgx.NamedArgs{
		"parent_id": parentID,
	}
	rows, err := r.db(ctx).Query(ctx, query, args)
	if err != nil {
		return nil, apperror.Internal("failed to get categories by parent", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
}