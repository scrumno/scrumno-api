package cancel_order

type Command struct {
	OrderID           string
	OrganizationID    string
	RemovalTypeID     string
	RemovalComment    string
	UserIDForWriteoff string
}
