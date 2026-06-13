package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/testutil"
	"github.com/google/uuid"
)

type mockRepo struct {
	createFunc                 func(ctx context.Context, company *model.Company) error
	getByIDFunc                func(ctx context.Context, id uuid.UUID) (*model.Company, error)
	getAllFunc                 func(ctx context.Context, params model.QueryParams) ([]model.Company, error)
	updateFunc                 func(ctx context.Context, company *model.Company) error
	deleteFunc                 func(ctx context.Context, id uuid.UUID) error
	getCompanyDependenciesFunc func(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error)
}

func (m *mockRepo) Create(ctx context.Context, company *model.Company) error {
	return m.createFunc(ctx, company)
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) GetAll(ctx context.Context, params model.QueryParams) ([]model.Company, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockRepo) Update(ctx context.Context, company *model.Company) error {
	return m.updateFunc(ctx, company)
}
func (m *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFunc(ctx, id)
}
func (m *mockRepo) GetCompanyDependencies(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error) {
	return m.getCompanyDependenciesFunc(ctx, id)
}

const emptyNameError = "company name is required"

func TestCreateCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, c *model.Company) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		company := testutil.CompanyMock()
		err := svc.CreateCompany(t.Context(), &company)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty name returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, c *model.Company) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		company := model.Company{CompanyName: ""}
		err := svc.CreateCompany(t.Context(), &company)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != emptyNameError {
			t.Errorf("expected '%s', got '%s'", emptyNameError, err.Error())
		}
	})
}

func TestGetCompanyByID(t *testing.T) {
	companyID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := testutil.CompanyMock()
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Company, error) {
				return &expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCompanyByID(t.Context(), companyID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.CompanyID != expected.CompanyID {
			t.Errorf("expected company ID %v, got %v", expected.CompanyID, got.CompanyID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Company, error) {
				return nil, errors.New("company not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetCompanyByID(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}

func TestGetAllCompanies(t *testing.T) {
	t.Run("success with results", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Company, error) {
				return []model.Company{testutil.CompanyMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllCompanies(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 company, got %d", len(got))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Company, error) {
				return []model.Company{}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllCompanies(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 companies, got %d", len(got))
		}
	})

	t.Run("negative limit returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllCompanies(t.Context(), model.QueryParams{Limit: -1, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be non-negative" {
			t.Errorf("expected 'limit must be non-negative', got '%s'", err.Error())
		}
	})

	t.Run("negative offset returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllCompanies(t.Context(), model.QueryParams{Limit: 10, Offset: -1})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "offset must be non-negative" {
			t.Errorf("expected 'offset must be non-negative', got '%s'", err.Error())
		}
	})

	t.Run("limit over max returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllCompanies(t.Context(), model.QueryParams{Limit: 101, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be less than or equal to 100" {
			t.Errorf("expected 'limit must be less than or equal to 100', got '%s'", err.Error())
		}
	})
}

func TestUpdateCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, c *model.Company) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		company := testutil.CompanyMock()
		err := svc.UpdateCompany(t.Context(), &company)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty name returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, c *model.Company) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		company := model.Company{CompanyName: ""}
		err := svc.UpdateCompany(t.Context(), &company)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != emptyNameError {
			t.Errorf("expected '%s', got '%s'", emptyNameError, err.Error())
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, c *model.Company) error {
				return errors.New("company not found")
			},
		}
		svc := service.NewService(mock)
		company := testutil.CompanyMock()
		err := svc.UpdateCompany(t.Context(), &company)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}

func TestGetCompanyDependencies(t *testing.T) {
	t.Run("success with dependencies", func(t *testing.T) {
		expected := model.CompanyDependency{ProductCount: 3, SupplierCount: 2}
		mock := &mockRepo{
			getCompanyDependenciesFunc: func(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error) {
				return expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCompanyDependencies(t.Context(), uuid.New())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ProductCount != expected.ProductCount {
			t.Errorf("expected product count %d, got %d", expected.ProductCount, got.ProductCount)
		}
		if got.SupplierCount != expected.SupplierCount {
			t.Errorf("expected supplier count %d, got %d", expected.SupplierCount, got.SupplierCount)
		}
	})

	t.Run("no dependencies", func(t *testing.T) {
		mock := &mockRepo{
			getCompanyDependenciesFunc: func(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error) {
				return model.CompanyDependency{}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCompanyDependencies(t.Context(), uuid.New())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ProductCount != 0 {
			t.Errorf("expected 0 products, got %d", got.ProductCount)
		}
		if got.SupplierCount != 0 {
			t.Errorf("expected 0 suppliers, got %d", got.SupplierCount)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getCompanyDependenciesFunc: func(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error) {
				return model.CompanyDependency{}, errors.New("company not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetCompanyDependencies(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}

func TestDeleteCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteCompany(t.Context(), uuid.New())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("company not found")
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteCompany(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}
