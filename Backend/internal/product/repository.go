package product

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetAll(ctx context.Context, params GetProductsQueryParams) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCompany(ctx context.Context, companyID uuid.UUID) ([]Product, error)
	GetByCategory(ctx context.Context, categoryID int) ([]Product, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, product *Product) error {
	query := `
		INSERT INTO "products" (product_name, product_description, diameter, width, company_id, price, category_id, location_id)
		VALUES (@product_name, @product_description, @diameter, @width, @company_id, @price, @category_id, @location_id)
		RETURNING product_id`
	args := pgx.NamedArgs{
		"product_name":        product.ProductName,
		"product_description": product.ProductDescription,
		"diameter":            product.Diameter,
		"width":               product.Width,
		"company_id":          product.CompanyID,
		"price":               product.Price,
		"category_id":         product.CategoryID,
		"location_id":         product.LocationID,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(&product.ProductID)

	if err != nil {
		return apperror.Internal("failed to create product", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	query := `
		SELECT product_id, product_name, product_description, diameter, width, company_id, price, category_id, location_id
		FROM "products"
		WHERE product_id = @product_id`

	var product Product
	err := r.db.QueryRow(ctx, query, pgx.NamedArgs{"product_id": id}).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductDescription,
		&product.Diameter,
		&product.Width,
		&product.CompanyID,
		&product.Price,
		&product.CategoryID,
		&product.LocationID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("product not found", err)
		}
		return nil, apperror.Internal("failed to get product by id", err)
	}
	return &product, nil
}

func (r *repository) GetAll(ctx context.Context, params GetProductsQueryParams) ([]Product, error) {
	query := `
		SELECT 
			i.stock,
			p.product_id,
			p.product_name, 
			p.product_description, 
			p.diameter, 
			p.width, 
			p.company_id, 
			p.price, 
			p.category_id, 
			p.location_id
		FROM "products" p 
		JOIN "inventories" i 
		on p.product_id = i.product_id 
		where 
			(@company_id::uuid IS NULL OR p.company_id = @company_id::uuid) 
			AND 
			(@category_id::int IS NULL OR p.category_id = @category_id::int)
		ORDER BY 
			p.product_name`

	rows, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"company_id":  params.CompanyID,
			"category_id": params.CategoryID,
		},
	)
	if err != nil {
		return nil, apperror.Internal("failed to get all products", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Product])

}

func (r *repository) Update(ctx context.Context, product *Product) error {
	query := `
		UPDATE "products"
		SET product_name = @product_name, product_description = @product_description, diameter = @diameter, width = @width, company_id = @company_id, price = @price, category_id = @category_id, location_id = @location_id
		WHERE product_id = @product_id`

	args := pgx.NamedArgs{
		"product_id":          product.ProductID,
		"product_name":        product.ProductName,
		"product_description": product.ProductDescription,
		"diameter":            product.Diameter,
		"width":               product.Width,
		"company_id":          product.CompanyID,
		"price":               product.Price,
		"category_id":         product.CategoryID,
		"location_id":         product.LocationID,
	}
	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to update product", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("product not found", nil)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "products" WHERE product_id = @product_id`

	args := pgx.NamedArgs{
		"product_id": id,
	}
	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to delete product", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("product not found", nil)
	}
	return nil
}

func (r *repository) GetByCompany(ctx context.Context, companyID uuid.UUID) ([]Product, error) {
	query := `
		SELECT product_id, product_name, product_description, diameter, width, company_id, price, category_id, location_id
		FROM "products"
		WHERE company_id = @company_id
		ORDER BY product_name`

	args := pgx.NamedArgs{
		"company_id": companyID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, apperror.Internal("failed to get products by company", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Product])
}

func (r *repository) GetByCategory(ctx context.Context, categoryID int) ([]Product, error) {
	query := `
		SELECT product_id, product_name, product_description, diameter, width, company_id, price, category_id, location_id
		FROM "products"
		WHERE category_id = @category_id
		ORDER BY product_name`

	args := pgx.NamedArgs{
		"category_id": categoryID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, apperror.Internal("failed to get products by category", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Product])
}
