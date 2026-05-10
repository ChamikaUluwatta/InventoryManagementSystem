package seed

import (
	"context"
	"log"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	companyModel "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	inventoryModel "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	locationModel "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	productModel "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type companyCreator interface {
	Create(ctx context.Context, company *companyModel.Company) error
}

type categoryCreator interface {
	Create(ctx context.Context, category *model.Category) error
}

type locationCreator interface {
	Create(ctx context.Context, location *locationModel.Location) error
}

type productCreator interface {
	Create(ctx context.Context, product *productModel.Product) error
}

type inventoryCreator interface {
	Create(ctx context.Context, inventory *inventoryModel.Inventory) error
}

func seedCompaniesFn(ctx context.Context, creator companyCreator, seeds []companySeed) ([]companyModel.Company, error) {
	var created []companyModel.Company
	for _, s := range seeds {
		company := companyModel.Company{CompanyName: s.Name, Description: s.Description}
		if err := creator.Create(ctx, &company); err != nil {
			log.Printf("Failed to create company %q: %v", s.Name, err)
			continue
		}
		created = append(created, company)
	}
	return created, nil
}

func seedCategoriesFn(ctx context.Context, creator categoryCreator, seeds []categorySeed) ([]model.Category, error) {
	var created []model.Category
	for _, s := range seeds {
		category := model.Category{CategoryName: s.Name}
		if err := creator.Create(ctx, &category); err != nil {
			log.Printf("Failed to create category %q: %v", s.Name, err)
			continue
		}
		created = append(created, category)
	}
	return created, nil
}

func seedLocationsFn(ctx context.Context, creator locationCreator, seeds []locationSeed) ([]locationModel.Location, error) {
	var created []locationModel.Location
	for _, s := range seeds {
		location := locationModel.Location{LocationID: s.LocationID}
		if err := creator.Create(ctx, &location); err != nil {
			log.Printf("Failed to create location %q: %v", s.LocationID, err)
			continue
		}
		created = append(created, location)
	}
	return created, nil
}

func seedProductsFn(ctx context.Context, creator productCreator, seeds []productSeed, companyIDs []uuid.UUID, categoryIDs []int) ([]productModel.Product, error) {
	if len(companyIDs) < 1 || len(categoryIDs) < 1 {
		return nil, nil
	}

	var created []productModel.Product
	for i, s := range seeds {
		product := productModel.Product{
			ProductName:        s.Name,
			ProductDescription: s.Description,
			Diameter:           decimal.NewFromFloat(s.Diameter),
			Width:              decimal.NewFromFloat(s.Width),
			CompanyID:          companyIDs[i%len(companyIDs)],
			Price:              decimal.NewFromFloat(s.Price),
			CategoryID:         categoryIDs[i%len(categoryIDs)],
		}
		if err := creator.Create(ctx, &product); err != nil {
			log.Printf("Failed to create product %q: %v", s.Name, err)
			continue
		}
		created = append(created, product)
	}
	return created, nil
}

func seedInventoriesFn(ctx context.Context, creator inventoryCreator, seeds []inventorySeed, productIDs []uuid.UUID, locationIDs []string) ([]inventoryModel.Inventory, error) {
	if len(productIDs) < 1 || len(locationIDs) < 1 {
		return nil, nil
	}

	var created []inventoryModel.Inventory
	for _, s := range seeds {
		inv := inventoryModel.Inventory{
			ProductID:  productIDs[s.ProductIndex%len(productIDs)],
			LocationID: locationIDs[s.LocationIndex%len(locationIDs)],
			Stock:      s.Stock,
		}
		if err := creator.Create(ctx, &inv); err != nil {
			log.Printf("Failed to create inventory: %v", err)
			continue
		}
		created = append(created, inv)
	}
	return created, nil
}
