package category

import (
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	base.BaseRepository[Category]
}

type categoryRepository struct {
	base.BaseRepository[Category]
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		BaseRepository: factory.NewGormRepository[Category](db),
	}
}
