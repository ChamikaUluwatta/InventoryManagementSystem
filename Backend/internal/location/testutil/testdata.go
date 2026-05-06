package testutil

import "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"

func LocationMock() model.Location {
	return model.Location{
		LocationID: "TEST-LOC-1",
	}
}

func LocationWithImageMock() model.Location {
	img := "https://example.com/image.png"
	return model.Location{
		LocationID: "TEST-LOC-2",
		Image:      &img,
	}
}
