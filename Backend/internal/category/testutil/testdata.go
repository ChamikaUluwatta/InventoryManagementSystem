package testutil

import "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"

func CategoryMock() model.Category {
	return model.Category{
		CategoryID:   1,
		CategoryName: "Test Category",
	}
}

func CategoryWithParentMock() model.Category {
	parentID := 1
	return model.Category{
		CategoryID:   2,
		CategoryName: "Subcategory",
		ParentID:     &parentID,
	}
}

func RootCategoryMock() model.Category {
	return model.Category{
		CategoryID:   3,
		CategoryName: "Root Category",
		ParentID:     nil,
	}
}
