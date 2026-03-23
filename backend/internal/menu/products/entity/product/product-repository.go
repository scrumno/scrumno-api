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
	base.BaseRepository[Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		BaseRepository: factory.NewGormRepository[Product](db),
	}
}
