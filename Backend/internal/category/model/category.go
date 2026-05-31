package model

import "github.com/google/uuid"

type Category struct {
	CategoryID   int       `db:"category_id"   json:"category_id"`
	CategoryName string    `db:"category_name" json:"category_name"`
	ParentID     *int      `db:"parent_id"     json:"parent_id"`
	TenantID     uuid.UUID `db:"tenant_id"     json:"-"`
}

type QueryParams struct {
	Limit  int
	Offset int
}