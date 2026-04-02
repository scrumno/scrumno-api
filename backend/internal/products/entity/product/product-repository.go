package product

import (
	"context"

	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type ProductRepository interface {
	base.BaseRepository[Product]

	FindByExternalID(ctx context.Context, externalID string) (*Product, error)
	GetByCategoryID(ctx context.Context, categoryID string, offset int, limit int) ([]Product, error)
}

type productRepository struct {
	*factory.GormRepository[Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		GormRepository: factory.NewGormRepository[Product](db),
	}
}

func (r *productRepository) FindByExternalID(ctx context.Context, externalID string) (*Product, error) {
	var entity Product
	err := r.DB.WithContext(ctx).Where("external_id = ?", externalID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID string, offset int, limit int) ([]Product, error) {
	var entities []Product
	err := r.DB.WithContext(ctx).Where("product_category_id = ?", categoryID).Offset(offset).Limit(limit).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
