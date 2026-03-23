package product

import (
	productAllergen "github.com/scrumno/scrumno-api/internal/menu/products/entity/product-allergen"
	customerTag "github.com/scrumno/scrumno-api/internal/menu/products/entity/product-customer-tag"
	productSize "github.com/scrumno/scrumno-api/internal/menu/products/entity/product-size"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type ProductRepository struct {
	SizeRepository        base.BaseRepository[productSize.ProductSize]
	AllergenRepository    base.BaseRepository[productAllergen.ProductAllergen]
	CustomerTagRepository base.BaseRepository[customerTag.ProductCustomerTag]
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		SizeRepository:        factory.NewGormRepository[productSize.ProductSize](db),
		AllergenRepository:    factory.NewGormRepository[productAllergen.ProductAllergen](db),
		CustomerTagRepository: factory.NewGormRepository[customerTag.ProductCustomerTag](db),
	}
}
