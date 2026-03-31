package category

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
	GetByParent(ctx context.Context, parentID *int) ([]Category, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, category *Category) error {
	query := `
		INSERT INTO "categories" (category_name, parent_id)
		VALUES (@category_name, @parent_id)
		RETURNING category_id`

	args := pgx.NamedArgs{
		"category_name": category.CategoryName,
		"parent_id":     category.ParentID,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(&category.CategoryID)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*Category, error) {
	query := `
		SELECT category_id, category_name, parent_id
		FROM "categories"
		WHERE category_id = @category_id`

	var category Category
	args := pgx.NamedArgs{
		"category_id": id,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(
		&category.CategoryID,
		&category.CategoryName,
		&category.ParentID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &category, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Category, error) {
	query := `
		SELECT category_id, category_name, parent_id
		FROM "categories"
		ORDER BY category_name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Category])
}

func (r *repository) Update(ctx context.Context, category *Category) error {
	query := `
		UPDATE "categories"
		SET category_name = @category_name, parent_id = @parent_id
		WHERE category_id = @category_id`
	args := pgx.NamedArgs{
		"category_name": category.CategoryName,
		"parent_id":     category.ParentID,
		"category_id":   category.CategoryID,
	}
	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "categories" WHERE category_id = @category_id`
	args := pgx.NamedArgs{
		"category_id": id,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (r *repository) GetByParent(ctx context.Context, parentID *int) ([]Category, error) {
	query := `
		SELECT category_id, category_name, parent_id
		FROM "categories"
		WHERE parent_id = @parent_id
		ORDER BY category_name`

	args := pgx.NamedArgs{
		"parent_id": parentID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by parent: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Category])
}
