package entity

type Customer struct {
	ID                                    string `json:"id,omitempty"`
	Name                                  string `json:"name,omitempty"`
	Surname                               string `json:"surname,omitempty"`
	Comment                               string `json:"comment,omitempty"`
	Birthdate                             string `json:"birthdate,omitempty"`
	Email                                 string `json:"email,omitempty"`
	ShouldReceivePromoActionsInfo         bool   `json:"shouldReceivePromoActionsInfo,omitempty"`
	ShouldReceiveOrderStatusNotifications bool   `json:"shouldReceiveOrderStatusNotifications,omitempty"`
	Gender                                string `json:"gender,omitempty"`
	Type                                  string `json:"type,omitempty"`
}

// Guests — count обязателен, если объект передан.
type Guests struct {
	Count               int   `json:"count"`
	SplitBetweenPersons *bool `json:"splitBetweenPersons,omitempty"`
}
