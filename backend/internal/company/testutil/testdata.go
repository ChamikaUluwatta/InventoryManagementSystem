package testutil

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/google/uuid"
)

func CompanyMock() model.Company {
	return model.Company{
		CompanyID:   uuid.New(),
		CompanyName: "Test Company",
	}
}
