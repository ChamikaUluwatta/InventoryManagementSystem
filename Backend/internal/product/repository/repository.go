package repository

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, product *model.CreateProductRequest) (model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error)
	GetAll(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, product *model.CreateProductRequest) (model.Product, error) {
	var productID uuid.UUID
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
	err := r.db.QueryRow(ctx, query, args).Scan(&productID)

	if err != nil {
		return model.Product{}, apperror.Internal("failed to create product", err)
	}
	return model.Product{ProductID: productID}, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
	query := `
		SELECT 
			p.product_id, 
			p.product_name, 
			p.product_description, 
			p.diameter, 
			p.width, 
			p.company_id, 
			p.price, 
			p.category_id, 
			p.location_id,
			COALESCE(i.stock, 0)
		FROM 
			"products" p
		LEFT JOIN 
			"inventories" i
		ON
			p.product_id = i.product_id
		WHERE 
			p.product_id = @product_id`

	var product model.GetProductById
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
		&product.Stock,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("product not found", err)
		}
		return nil, apperror.Internal("failed to get product by id", err)
	}
	return &product, nil
}

func (r *repository) GetAll(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
	query := `
		SELECT 
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
		where 
			(@company_id::uuid IS NULL OR p.company_id = @company_id::uuid) 
			AND 
			(@category_id::int IS NULL OR p.category_id = @category_id::int)
		ORDER BY 
			p.product_name
		LIMIT @limit OFFSET @offset`

	rows, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"company_id":  params.CompanyID,
			"category_id": params.CategoryID,
			"limit":       params.Limit,
			"offset":      params.Offset,
		},
	)
	if err != nil {
		return nil, apperror.Internal("failed to get all products", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Product])

}

func (r *repository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE "products"
		SET 
			product_name = @product_name, 
			product_description = @product_description, 
			diameter = @diameter, 
			width = @width, 
			company_id = @company_id, 
			price = @price, 
			category_id = @category_id, 
			location_id = @location_id
		WHERE 
			product_id = @product_id`

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
