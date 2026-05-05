package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	companyCreator   companyCreator
	categoryCreator  categoryCreator
	locationCreator  locationCreator
	productCreator   productCreator
	inventoryCreator inventoryCreator
	db               *pgxpool.Pool
}

func NewService(
	companyCreator companyCreator,
	categoryCreator categoryCreator,
	locationCreator locationCreator,
	productCreator productCreator,
	inventoryCreator inventoryCreator,
	db *pgxpool.Pool,
) *Service {
	return &Service{
		companyCreator:   companyCreator,
		categoryCreator:  categoryCreator,
		locationCreator:  locationCreator,
		productCreator:   productCreator,
		inventoryCreator: inventoryCreator,
		db:               db,
	}
}

type SeedResult struct {
	Companies   int `json:"companies_created"`
	Categories  int `json:"categories_created"`
	Locations   int `json:"locations_created"`
	Products    int `json:"products_created"`
	Inventories int `json:"inventories_created"`
}

type SeededIDs struct {
	CompanyIDs  []uuid.UUID `json:"company_ids"`
	CategoryIDs []int       `json:"category_ids"`
	LocationIDs []string    `json:"location_ids"`
	ProductIDs  []uuid.UUID `json:"product_ids"`
}

func (s *Service) Seed(ctx context.Context) (*SeedResult, *SeededIDs, error) {
	if err := s.clearTables(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to clear tables: %w", err)
	}

	result := &SeedResult{}
	ids := &SeededIDs{}

	companies, err := seedCompaniesFn(ctx, s.companyCreator, defaultCompanies)
	if err != nil {
		return nil, nil, err
	}
	result.Companies = len(companies)
	for _, c := range companies {
		ids.CompanyIDs = append(ids.CompanyIDs, c.CompanyID)
	}

	categories, err := seedCategoriesFn(ctx, s.categoryCreator, defaultCategories)
	if err != nil {
		return nil, nil, err
	}
	result.Categories = len(categories)
	for _, c := range categories {
		ids.CategoryIDs = append(ids.CategoryIDs, c.CategoryID)
	}

	locations, err := seedLocationsFn(ctx, s.locationCreator, defaultLocations)
	if err != nil {
		return nil, nil, err
	}
	result.Locations = len(locations)
	for _, l := range locations {
		ids.LocationIDs = append(ids.LocationIDs, l.LocationID)
	}

	products, err := seedProductsFn(ctx, s.productCreator, defaultProducts, ids.CompanyIDs, ids.CategoryIDs)
	if err != nil {
		return nil, nil, err
	}
	result.Products = len(products)
	for _, p := range products {
		ids.ProductIDs = append(ids.ProductIDs, p.ProductID)
	}

	inventories, err := seedInventoriesFn(ctx, s.inventoryCreator, defaultInventories, ids.ProductIDs, ids.LocationIDs)
	if err != nil {
		return nil, nil, err
	}
	result.Inventories = len(inventories)

	return result, ids, nil
}

func (s *Service) clearTables(ctx context.Context) error {
	_, err := s.db.Exec(ctx, "TRUNCATE TABLE inventories, products, locations, categories, companies RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}
	log.Println("All tables cleared successfully")
	return nil
}
