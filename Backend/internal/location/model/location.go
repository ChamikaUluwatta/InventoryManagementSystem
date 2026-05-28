package model

import "github.com/google/uuid"

type Location struct {
	LocationID string    `db:"location_id" json:"location_id"`
	Image      *string   `db:"image"       json:"image"`
	TenantID   uuid.UUID `db:"tenant_id"   json:"-"`
}

type QueryParams struct {
	Limit  int
	Offset int
}
