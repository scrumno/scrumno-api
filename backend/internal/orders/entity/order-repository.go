package entity

func create() *DeliveryOrder {
	return &DeliveryOrder{
		Phone: "",
		Items: []OrderItem{},
	}
}
