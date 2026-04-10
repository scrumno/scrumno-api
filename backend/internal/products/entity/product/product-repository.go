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
	FindCookingTimeProductTableByExternalID(ctx context.Context, externalID string) (*CookingTimeProductTable, error)
	UpdateCookingTimeProductTable(ctx context.Context, entity *CookingTimeProductTable) error
	GetCookingTimeProductTableByProductIDs(ctx context.Context, productIDs []string) ([]CookingTimeProductTable, error)
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

func (r *productRepository) GetCookingTimeProductTable(ctx context.Context, productID string) (*CookingTimeProductTable, error) {
	var entity CookingTimeProductTable
	err := r.DB.WithContext(ctx).Where("product_id = ?", productID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) UpdateCookingTimeProductTable(ctx context.Context, entity *CookingTimeProductTable) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *productRepository) GetCookingTimeProductTableByProductIDs(ctx context.Context, productIDs []string) ([]CookingTimeProductTable, error) {
	var entities []CookingTimeProductTable
	err := r.DB.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *productRepository) FindCookingTimeProductTableByExternalID(ctx context.Context, externalID string) (*CookingTimeProductTable, error) {
	var entity CookingTimeProductTable
	err := r.DB.WithContext(ctx).Where("external_id = ?", externalID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
