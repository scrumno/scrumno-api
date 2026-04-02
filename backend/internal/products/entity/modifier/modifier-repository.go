package modifier

import (
	"context"

	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"gorm.io/gorm"
)

type ModifierRepository interface {
	SaveProductModifier(ctx context.Context, modifier *ProductModifier) error
	SaveProductChildModifier(ctx context.Context, modifier *ProductChildModifier) error
	SaveProductModifierGroup(ctx context.Context, group *ProductModifierGroup) error
}

type modifierRepository struct {
	ProductModifierRepo      *factory.GormRepository[ProductModifier]
	ProductChildModifierRepo *factory.GormRepository[ProductChildModifier]
	ProductModifierGroupRepo *factory.GormRepository[ProductModifierGroup]
}

func NewModifierRepository(db *gorm.DB) *modifierRepository {
	return &modifierRepository{
		ProductModifierRepo:      factory.NewGormRepository[ProductModifier](db),
		ProductChildModifierRepo: factory.NewGormRepository[ProductChildModifier](db),
		ProductModifierGroupRepo: factory.NewGormRepository[ProductModifierGroup](db),
	}
}

func (r *modifierRepository) SaveProductModifier(ctx context.Context, modifier *ProductModifier) error {
	_, err := r.ProductModifierRepo.Save(ctx, modifier)
	return err
}

func (r *modifierRepository) SaveProductChildModifier(ctx context.Context, modifier *ProductChildModifier) error {
	_, err := r.ProductChildModifierRepo.Save(ctx, modifier)
	return err
}

func (r *modifierRepository) SaveProductModifierGroup(ctx context.Context, group *ProductModifierGroup) error {
	_, err := r.ProductModifierGroupRepo.Save(ctx, group)
	return err
}
