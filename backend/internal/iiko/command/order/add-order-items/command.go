package add_order_items

type Command struct {
	OrderID        string
	OrganizationID string
	Items          []Item
}

type Item struct {
	ProductID string
	Amount    float64
	Price     float64
	Comment   string
}
