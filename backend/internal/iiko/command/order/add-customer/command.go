package add_customer

type Command struct {
	OrganizationID string
	OrderID        string
	Customer       Customer
}

type Customer struct {
	ID                                    string
	Name                                  string
	Surname                               string
	Comment                               string
	Birthdate                             *string
	Email                                 string
	ShouldReceiveOrderStatusNotifications *bool
	Gender                                string
	Phone                                 string
}
