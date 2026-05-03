package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Service struct {
	companyRepo   company.Repository
	categoryRepo  category.Repository
	locationRepo  location.Repository
	productRepo   product.Repository
	inventoryRepo inventory.Repository
	db            *pgxpool.Pool
}

func NewService(
	companyRepo company.Repository,
	categoryRepo category.Repository,
	locationRepo location.Repository,
	productRepo product.Repository,
	inventoryRepo inventory.Repository,
	db *pgxpool.Pool,
) *Service {
	return &Service{
		companyRepo:   companyRepo,
		categoryRepo:  categoryRepo,
		locationRepo:  locationRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		db:            db,
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

	companies, err := s.seedCompanies(ctx)
	if err != nil {
		return nil, nil, err
	}
	result.Companies = len(companies)
	for _, c := range companies {
		ids.CompanyIDs = append(ids.CompanyIDs, c.CompanyID)
	}

	categories, err := s.seedCategories(ctx)
	if err != nil {
		return nil, nil, err
	}
	result.Categories = len(categories)
	for _, c := range categories {
		ids.CategoryIDs = append(ids.CategoryIDs, c.CategoryID)
	}

	locations, err := s.seedLocations(ctx)
	if err != nil {
		return nil, nil, err
	}
	result.Locations = len(locations)
	for _, l := range locations {
		ids.LocationIDs = append(ids.LocationIDs, l.LocationID)
	}

	products, err := s.seedProducts(ctx, ids.CompanyIDs, ids.CategoryIDs)
	if err != nil {
		return nil, nil, err
	}
	result.Products = len(products)
	for _, p := range products {
		ids.ProductIDs = append(ids.ProductIDs, p.ProductID)
	}

	inventories, err := s.seedInventories(ctx, ids.ProductIDs, ids.LocationIDs)
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

func (s *Service) seedCompanies(ctx context.Context) ([]company.Company, error) {
	companies := []company.Company{
		{CompanyName: "Acme Corp"},
		{CompanyName: "Tech Solutions"},
		{CompanyName: "Global Industries"},
	}

	var created []company.Company
	for i := range companies {
		if err := s.companyRepo.Create(ctx, &companies[i]); err != nil {
			log.Printf("Failed to create company: %v", err)
		} else {
			created = append(created, companies[i])
		}
	}
	return created, nil
}

func (s *Service) seedCategories(ctx context.Context) ([]category.Category, error) {
	categories := []category.Category{
		{CategoryName: "Electronics"},
		{CategoryName: "Hardware"},
		{CategoryName: "Tools"},
	}

	var created []category.Category
	for i := range categories {
		if err := s.categoryRepo.Create(ctx, &categories[i]); err != nil {
			log.Printf("Failed to create category: %v", err)
		} else {
			created = append(created, categories[i])
		}
	}
	return created, nil
}

func (s *Service) seedLocations(ctx context.Context) ([]location.Location, error) {
	locations := []location.Location{
		{LocationID: "LOC-001", Image: nil},
		{LocationID: "LOC-002", Image: nil},
		{LocationID: "WAREHOUSE-A", Image: nil},
	}

	var created []location.Location
	for i := range locations {
		if err := s.locationRepo.Create(ctx, &locations[i]); err != nil {
			log.Printf("Failed to create location: %v", err)
		} else {
			created = append(created, locations[i])
		}
	}
	return created, nil
}

func (s *Service) seedProducts(ctx context.Context, companyIDs []uuid.UUID, categoryIDs []int) ([]product.Product, error) {
	if len(companyIDs) < 1 || len(categoryIDs) < 1 {
		return nil, nil
	}

<<<<<<< Updated upstream
	products := []product.Product{
=======
	products := []productModel.CreateProductRequest{
>>>>>>> Stashed changes
		{
			ProductName: "Widget A",
			Diameter:    decimal.NewFromFloat(10.5),
			Width:       decimal.NewFromFloat(5.0),
			CompanyID:   companyIDs[0],
			Price:       decimal.NewFromFloat(99.99),
			CategoryID:  categoryIDs[0],
		},
		{
			ProductName: "Gadget B",
			Diameter:    decimal.NewFromFloat(8.0),
			Width:       decimal.NewFromFloat(3.5),
			CompanyID:   companyIDs[1%len(companyIDs)],
			Price:       decimal.NewFromFloat(149.99),
			CategoryID:  categoryIDs[1%len(categoryIDs)],
		},
		{
			ProductName: "Tool C",
			Diameter:    decimal.NewFromFloat(12.0),
			Width:       decimal.NewFromFloat(6.0),
			CompanyID:   companyIDs[2%len(companyIDs)],
			Price:       decimal.NewFromFloat(79.99),
			CategoryID:  categoryIDs[2%len(categoryIDs)],
		},
	}

	var created []product.Product
	for i := range products {
		if product, err := s.productRepo.Create(ctx, &products[i]); err != nil {
			log.Printf("Failed to create product: %v", err)
		} else {
			created = append(created, product)
		}
	}
	return created, nil
}

func (s *Service) seedInventories(ctx context.Context, productIDs []uuid.UUID, locationIDs []string) ([]inventory.Inventory, error) {
	if len(productIDs) < 1 || len(locationIDs) < 1 {
		return nil, nil
	}

	inventories := []inventory.Inventory{
		{
			ProductID:  productIDs[0],
			LocationID: locationIDs[0],
			Stock:      100,
		},
		{
			ProductID:  productIDs[1%len(productIDs)],
			LocationID: locationIDs[1%len(locationIDs)],
			Stock:      50,
		},
		{
			ProductID:  productIDs[2%len(productIDs)],
			LocationID: locationIDs[2%len(locationIDs)],
			Stock:      200,
		},
	}

	var created []inventory.Inventory
	for i := range inventories {
		if err := s.inventoryRepo.Create(ctx, &inventories[i]); err != nil {
			log.Printf("Failed to create inventory: %v", err)
		} else {
			created = append(created, inventories[i])
		}
	}
	return created, nil
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Seed(w http.ResponseWriter, r *http.Request) {
	result, ids, err := h.service.Seed(r.Context())
	if err != nil {
		apperror.HandleError(w, apperror.Internal("seed failed", err))
		return
	}

	response := map[string]interface{}{
		"message": "Seed completed successfully",
		"result":  result,
		"ids":     ids,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /seed", h.Seed)
}
