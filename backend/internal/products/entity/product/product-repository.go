package product

import (
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type ProductRepository interface {
	base.BaseRepository[Product]
}

type productRepository struct {
	*factory.GormRepository[Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		GormRepository: factory.NewGormRepository[Product](db),
	}
}
