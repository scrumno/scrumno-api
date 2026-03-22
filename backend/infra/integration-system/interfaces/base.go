package interfaces

type OrderBuilder interface {
	BuildBody(data *any) *any
}

type OrderService interface {
	Create(order *any) error
}
