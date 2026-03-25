package interfaces

type OrderProvider interface {
	Create(order any) (any, error)
}

type OrderBuilder interface {
	BuildBody(data any) any
}
