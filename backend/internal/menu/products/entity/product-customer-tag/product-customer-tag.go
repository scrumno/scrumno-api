package product_customer_tag

import "time"

type ProductCustomerTag struct {
	ID uint `gorm:"primaryKey"`

	ProductID     uint `gorm:"index:idx_product_customer_tag,unique;not null"`
	CustomerTagID uint `gorm:"index:idx_product_customer_tag,unique;not null"`

	CreatedAt time.Time
}

