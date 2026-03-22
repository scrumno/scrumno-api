package create_order

type Command struct {
	OrganizationID   string
	TerminalGroupID  string
	CustomerName     string
	CustomerPhone    string
	Comment          string
	CompleteBefore   *string
	OrderServiceType *string
	Items            []Item
}

type Item struct {
	ProductID string
	Amount    float64
	Price     float64
	Comment   string
}
