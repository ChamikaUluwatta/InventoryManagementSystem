package model

type Category struct {
	CategoryID   int    `db:"category_id"   json:"category_id"`
	CategoryName string `db:"category_name" json:"category_name"`
	ParentID     *int   `db:"parent_id"     json:"parent_id"`
}

type QueryParams struct {
	Limit  int
	Offset int
}