package company

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, company *Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*Company, error)
	GetAll(ctx context.Context) ([]Company, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, company *Company) error {
	query := `
		INSERT INTO "companies" (company_name)
		VALUES (@company_name)
		RETURNING company_id`

	args := pgx.NamedArgs{
		"company_name": company.CompanyName,
	}
	err := r.db.QueryRow(ctx, query, args).Scan(&company.CompanyID)

	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Company, error) {
	query := `
		SELECT company_id, company_name
		FROM "companies"
		WHERE company_id = @company_id`

	args := pgx.NamedArgs{
		"company_id": id,
	}
	var company Company
	err := r.db.QueryRow(ctx, query, args).Scan(
		&company.CompanyID,
		&company.CompanyName,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get company by id: %w", err)
	}
	return &company, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Company, error) {
	query := `
		SELECT company_id, company_name
		FROM "companies"
		ORDER BY company_name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all companies: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Company])
}

func (r *repository) Update(ctx context.Context, company *Company) error {
	query := `
		UPDATE "companies"
		SET company_name = @company_name
		WHERE company_id = @company_id`

	args := pgx.NamedArgs{
		"company_name": company.CompanyName,
		"company_id":   company.CompanyID,
	}
	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "companies" WHERE company_id = @company_id`
	args := pgx.NamedArgs{
		"company_id": id,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}
