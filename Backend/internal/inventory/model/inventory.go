package model

import "github.com/google/uuid"

type Inventory struct {
	InventoryID int       `db:"inventory_id" json:"inventory_id"`
	ProductID   uuid.UUID `db:"product_id"   json:"product_id"`
	LocationID  string    `db:"location_id"  json:"location_id"`
	Stock       int       `db:"stock"        json:"stock"`
}