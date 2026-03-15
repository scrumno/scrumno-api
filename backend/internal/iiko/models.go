package iiko

import "time"

// PickupOrder описывает упрощённую модель заказа самовывоза,
// с которой работает наш бэкенд (вход от клиента/POSTMAN).
type PickupOrder struct {
	Customer Customer     `json:"customer"`
	Items    []OrderItem  `json:"items"`
	Comment  string       `json:"comment,omitempty"`
	PickupAt *time.Time   `json:"pickupAt,omitempty"`
}

type Customer struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type OrderItem struct {
	ProductID string  `json:"productId"`
	Quantity  float64 `json:"quantity"`
}

// CreateOrderResult содержит минимальную информацию о результате запроса в iiko.
type CreateOrderResult struct {
	StatusCode int
	Body       []byte
}

