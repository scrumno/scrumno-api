package entity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	except "github.com/scrumno/scrumno-api/shared/exception/cart"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type CartRepository interface {
	base.BaseRepository[Cart]
	AddProductToCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int, price float64) error
	GetCartByUserId(ctx context.Context, userID uuid.UUID) (*Cart, error)
	UpdateCartProduct(ctx context.Context, productID uuid.UUID, quantity int, price float64, cart *Cart) (bool, error)
}

type cartGormRepository struct {
	*factory.GormRepository[Cart]
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartGormRepository{
		GormRepository: factory.NewGormRepository[Cart](db),
	}
}

func (r *cartGormRepository) AddProductToCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int, price float64) error {

	cart, err := r.GetCartByUserId(ctx, userID)
	if err != nil {
		return err
	}

	updated, err := r.UpdateCartProduct(ctx, productID, quantity, price, cart)
	if err != nil {
		return err
	}

	if updated {
		return nil
	}

	newProductItem := CartItem{
		ID:        uuid.New(),
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.DB.WithContext(ctx).Create(&newProductItem).Error; err != nil {
		return except.ErrCartAddProduct
	}
	return nil
}

func (r *cartGormRepository) UpdateCartProduct(ctx context.Context, productID uuid.UUID, quantity int, price float64, cart *Cart) (bool, error) {
	var existingItem *CartItem
	for i, item := range cart.Items {
		if item.ProductID == productID {
			existingItem = &cart.Items[i]
			break
		}
	}

	if existingItem != nil {
		existingItem.Quantity += quantity
		existingItem.UnitPrice = price
		existingItem.UpdatedAt = time.Now()
		if err := r.DB.WithContext(ctx).Save(existingItem).Error; err != nil {
			return false, except.ErrCartUpdated
		}
		return true, nil
	}

	return false, nil
}

func (r *cartGormRepository) GetCartByUserId(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	var cart Cart

	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCart := NewCart(userID)
			if err := r.Create(ctx, newCart).Error; err != nil {
				return nil, except.ErrCartCreate
			}
			return newCart, nil
		}

		return nil, except.ErrCartFind
	}

	return &cart, nil
}
