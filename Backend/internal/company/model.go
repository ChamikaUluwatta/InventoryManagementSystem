package company

import "github.com/google/uuid"

type Company struct {
	CompanyID   uuid.UUID `db:"company_id"   json:"company_id"`
	CompanyName string    `db:"company_name" json:"company_name"`
}
