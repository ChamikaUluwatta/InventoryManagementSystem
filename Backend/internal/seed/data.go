package seed

type companySeed struct {
	Name string
}

type categorySeed struct {
	Name string
}

type locationSeed struct {
	LocationID string
}

type productSeed struct {
	Name        string
	Description string
	Diameter    float64
	Width       float64
	Price       float64
}

type inventorySeed struct {
	ProductIndex  int
	LocationIndex int
	Stock         int
}

var defaultCompanies = []companySeed{
	{Name: "Acme Corp"},
	{Name: "Tech Solutions"},
	{Name: "Global Industries"},
}

var defaultCategories = []categorySeed{
	{Name: "Electronics"},
	{Name: "Hardware"},
	{Name: "Tools"},
}

var defaultLocations = []locationSeed{
	{LocationID: "LOC-001"},
	{LocationID: "LOC-002"},
	{LocationID: "WAREHOUSE-A"},
}

var defaultProducts = []productSeed{
	{Name: "Widget A", Description: "A standard widget", Diameter: 10.5, Width: 5.0, Price: 99.99},
	{Name: "Gadget B", Description: "A fancy gadget", Diameter: 8.0, Width: 3.5, Price: 149.99},
	{Name: "Tool C", Description: "A handy tool", Diameter: 12.0, Width: 6.0, Price: 79.99},
}

var defaultInventories = []inventorySeed{
	{ProductIndex: 0, LocationIndex: 0, Stock: 100},
	{ProductIndex: 1, LocationIndex: 1, Stock: 50},
	{ProductIndex: 2, LocationIndex: 2, Stock: 200},
}
