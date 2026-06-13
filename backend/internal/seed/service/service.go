package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	companyCreator   CompanyCreator
	categoryCreator  CategoryCreator
	locationCreator  LocationCreator
	productCreator   ProductCreator
	inventoryCreator InventoryCreator
	pool             *pgxpool.Pool
}

func NewService(
	companyCreator CompanyCreator,
	categoryCreator CategoryCreator,
	locationCreator LocationCreator,
	productCreator ProductCreator,
	inventoryCreator InventoryCreator,
	pool *pgxpool.Pool,
) *Service {
	return &Service{
		companyCreator:   companyCreator,
		categoryCreator:  categoryCreator,
		locationCreator:  locationCreator,
		productCreator:   productCreator,
		inventoryCreator: inventoryCreator,
		pool:             pool,
	}
}

func (s *Service) db(ctx context.Context) database.Querier {
	if tx, ok := database.GetTx(ctx); ok {
		return tx
	}
	return s.pool
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

func (s *Service) Seed(ctx context.Context, tenantID uuid.UUID) (*SeedResult, *SeededIDs, error) {

	result := &SeedResult{}
	ids := &SeededIDs{}

	companies, err := seedCompaniesFn(ctx, s.companyCreator, defaultCompanies, tenantID)
	if err != nil {
		return nil, nil, err
	}
	result.Companies = len(companies)
	for _, c := range companies {
		ids.CompanyIDs = append(ids.CompanyIDs, c.CompanyID)
	}

	categories, err := seedCategoriesFn(ctx, s.categoryCreator, defaultCategories, tenantID)
	if err != nil {
		return nil, nil, err
	}
	result.Categories = len(categories)
	for _, c := range categories {
		ids.CategoryIDs = append(ids.CategoryIDs, c.CategoryID)
	}

	locations, err := seedLocationsFn(ctx, s.locationCreator, defaultLocations, tenantID)
	if err != nil {
		return nil, nil, err
	}
	result.Locations = len(locations)
	for _, l := range locations {
		ids.LocationIDs = append(ids.LocationIDs, l.LocationID)
	}

	products, err := seedProductsFn(ctx, s.productCreator, defaultProducts, ids.CompanyIDs, ids.CategoryIDs, ids.LocationIDs, tenantID)
	if err != nil {
		return nil, nil, err
	}
	result.Products = len(products)
	for _, p := range products {
		ids.ProductIDs = append(ids.ProductIDs, p.ProductID)
	}

	inventories, err := seedInventoriesFn(ctx, s.inventoryCreator, defaultInventories, ids.ProductIDs, ids.LocationIDs, tenantID)
	if err != nil {
		return nil, nil, err
	}
	result.Inventories = len(inventories)

	return result, ids, nil
}
