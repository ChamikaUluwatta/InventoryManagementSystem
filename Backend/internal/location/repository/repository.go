package repository

import (
	"context"
	"errors"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, location *model.Location) error
	GetAll(ctx context.Context, params model.QueryParams) ([]model.Location, error)
	GetByID(ctx context.Context, id string) (*model.Location, error)
	Update(ctx context.Context, location *model.Location) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, location *model.Location) error {
	query := `INSERT INTO "locations" (location_id,image) VALUES (@location_id,@image)`
	args := pgx.NamedArgs{
		"location_id": location.LocationID,
		"image":       location.Image,
	}
	_, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return apperror.Internal("failed to create location", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "locations" WHERE location_id = @location_id`

	args := pgx.NamedArgs{
		"location_id": id,
	}

	result, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return apperror.Internal("failed to delete location", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("location not found", nil)
	}
	return nil
}

func (r *repository) GetAll(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
	query := `SELECT location_id,image FROM "locations" ORDER BY location_id LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	rows, err := r.db.Query(ctx, query, args)

	if err != nil {
		return nil, apperror.Internal("failed to get all locations", err)
	}
	defer rows.Close()
	locations, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Location])

	if err != nil {
		return nil, apperror.Internal("failed to collect location rows", err)
	}

	return locations, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*model.Location, error) {
	query := `SELECT location_id,image FROM "locations" WHERE location_id = @location_id`
	args := pgx.NamedArgs{
		"location_id": id,
	}
	var location model.Location
	err := r.db.QueryRow(ctx, query, args).Scan(&location.LocationID, &location.Image)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("location not found", err)
		}
		return nil, apperror.Internal("failed to get location by id", err)
	}

	return &location, nil
}

func (r *repository) Update(ctx context.Context, location *model.Location) error {
	query := `
		UPDATE "locations"
		SET image = @image
		WHERE location_id = @location_id`
	args := pgx.NamedArgs{
		"location_id": location.LocationID,
		"image":       location.Image,
	}
	result, err := r.db.Exec(ctx, query, args)

	if err != nil {
		return apperror.Internal("failed to update location", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NotFound("location not found", nil)
	}
	return nil
}
