package create_order

type CartLineItem struct {
	ProductID string
	Quantity  float64
	Price     float64
	Comment   string
}

type Command struct {
	CustomerPhone    string
	CustomerFullName string
	CartItems        []CartLineItem
	SourceKey        string
	OrderComment     *string
}
