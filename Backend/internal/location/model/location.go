package model

type Location struct {
	LocationID string  `db:"location_id" json:"location_id"`
	Image      *string `db:"image"       json:"image"`
}