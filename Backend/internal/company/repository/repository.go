package repository

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, company *model.Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Company, error)
	GetAll(ctx context.Context, params model.QueryParams) ([]model.Company, error)
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetCompanyDependencies(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, company *model.Company) error {
	query := `
		INSERT INTO "companies" (company_name, description)
		VALUES (@company_name, @description)
		RETURNING company_id`

	args := pgx.NamedArgs{
		"company_name": company.CompanyName,
		"description":  company.Description,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(&company.CompanyID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperror.Conflict("company with the same name already exists", err)
		}
		return apperror.Internal("failed to create company", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	query := `
		SELECT company_id, company_name, description
		FROM "companies"
		WHERE company_id = @company_id`

	args := pgx.NamedArgs{
		"company_id": id,
	}
	var company model.Company
	err := r.db.QueryRow(ctx, query, args).Scan(
		&company.CompanyID,
		&company.CompanyName,
		&company.Description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("company not found", err)
		}
		return nil, apperror.Internal("failed to get company by id", err)
	}
	return &company, nil
}

func (r *repository) GetAll(ctx context.Context, params model.QueryParams) ([]model.Company, error) {
	query := `
		SELECT company_id, company_name, description
		FROM "companies"
		ORDER BY company_name
		LIMIT @limit OFFSET @offset`

	args := pgx.NamedArgs{
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, apperror.Internal("failed to get all companies", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Company])
}

func (r *repository) Update(ctx context.Context, company *model.Company) error {
	query := `
		UPDATE "companies"
		SET company_name = @company_name, description = @description
		WHERE company_id = @company_id`

	args := pgx.NamedArgs{
		"company_name": company.CompanyName,
		"description":  company.Description,
		"company_id":   company.CompanyID,
	}
	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to update company", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("company not found", nil)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "companies" WHERE company_id = @company_id`
	args := pgx.NamedArgs{
		"company_id": id,
	}

	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.NotFound("company not found", err)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return apperror.Conflict("cannot delete company with existing dependencies", err)
		}
		return apperror.Internal("failed to delete company", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("company not found", nil)
	}
	return nil
}

func (r *repository) GetCompanyDependencies(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM products WHERE company_id = @company_id) AS product_count,
			(SELECT COUNT(*) FROM supplier_returns WHERE company_id = @company_id) AS supplier_count`

	args := pgx.NamedArgs{
		"company_id": id,
	}
	var dep model.CompanyDependency
	err := r.db.QueryRow(ctx, query, args).Scan(&dep.ProductCount, &dep.SupplierCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.CompanyDependency{}, apperror.NotFound("company not found", err)
		}
		return model.CompanyDependency{}, apperror.Internal("failed to get company dependencies", err)
	}
	return dep, nil
}
