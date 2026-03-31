package location

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, location *Location) error
	GetAll(ctx context.Context) ([]Location, error)
	GetById(ctx context.Context, id string) (*Location, error)
	Update(ctx context.Context, location *Location) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, location *Location) error {
	query := `INSERT INTO "locations" (location_id,image) VALUES (@location_id,@image)`
	args := pgx.NamedArgs{
		"location_id": location.LocationID,
		"image":       location.Image,
	}
	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "locations" WHERE location_id = @location_id`

	args := pgx.NamedArgs{
		"location_id": id,
	}

	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}
	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]Location, error) {
	query := `SELECT location_id,image FROM "locations" ORDER BY location_id`
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("unable to query rows: %w", err)
	}
	defer rows.Close()
	Location, err := pgx.CollectRows(rows, pgx.RowToStructByName[Location])

	if err != nil {
		return nil, fmt.Errorf("unable to collect rows: %w", err)
	}

	return Location, nil
}

func (r *repository) GetById(ctx context.Context, id string) (*Location, error) {
	query := `SELECT location_id,image FROM "locations" WHERE location_id = @location_id`
	args := pgx.NamedArgs{
		"location_id": id,
	}
	var location Location
	err := r.db.QueryRow(ctx, query, args).Scan(&location.LocationID, &location.Image)

	if err != nil {
		return nil, fmt.Errorf("Invalid Location id %w", err)
	}

	return &location, nil
}

func (r *repository) Update(ctx context.Context, location *Location) error {
	query := `
		UPDATE "locations"
		SET image = @image
		WHERE location_id = @location_id`
	args := pgx.NamedArgs{
		"location_id": location.LocationID,
		"image":       location.Image,
	}
	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}
	return nil
}
