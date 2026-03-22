package entity

import (
	"context"

	"github.com/scrumno/scrumno-api/shared/base"
	"github.com/scrumno/scrumno-api/shared/factory"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	base.BaseRepository[User]
	FindByPhone(ctx context.Context, phone string) (*User, error)
	UpdateFieldsByPhone(ctx context.Context, phone string, fields map[string]any) error
}

type registrationGormRepository struct {
	*factory.GormRepository[User]
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationGormRepository{
		GormRepository: factory.NewGormRepository[User](db),
	}
}

func (r *registrationGormRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var u User
	err := r.DB.WithContext(ctx).
		Where("phone = ?", phone).
		First(&u).Error

	if err != nil {
		return nil, nil
	}

	return &u, nil
}

func (r *registrationGormRepository) UpdateFieldsByPhone(ctx context.Context, phone string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	return r.DB.WithContext(ctx).
		Model(&User{}).
		Where("phone = ?", phone).
		Updates(fields).Error
}
