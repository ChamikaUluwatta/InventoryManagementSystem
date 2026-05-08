package testutil

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/google/uuid"
)

func InventoryMock() model.Inventory {
	return model.Inventory{
		InventoryID: 1,
		ProductID:   uuid.New(),
		LocationID:  "TEST-LOC-1",
		Stock:       100,
	}
}

func CreateInventoryRequestMock() model.Inventory {
	return model.Inventory{
		ProductID:  uuid.New(),
		LocationID: "TEST-LOC-1",
		Stock:      50,
	}
}
