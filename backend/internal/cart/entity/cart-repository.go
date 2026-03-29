package entity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	except "github.com/scrumno/scrumno-api/shared/exception/cart"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

const (
	maxCartItemAmount = 999.999
	qtyEpsilon        = 1e-6
)

type CartRepository interface {
	base.BaseRepository[Cart]
	AddProductToCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity float64, unitPrice float64) error
	GetCartByUserId(ctx context.Context, userID uuid.UUID) (*Cart, error)
	UpdateCartProduct(ctx context.Context, productID uuid.UUID, quantityDelta float64, unitPrice float64, cart *Cart) (bool, error)
	RemoveProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type cartGormRepository struct {
	*factory.GormRepository[Cart]
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartGormRepository{
		GormRepository: factory.NewGormRepository[Cart](db),
	}
}

func validateCartItemAmount(q float64) error {
	if q <= 0 {
		return except.ErrCartInvalidQuantity
	}
	if q > maxCartItemAmount+qtyEpsilon {
		return except.ErrCartInvalidQuantity
	}
	return nil
}

func (r *cartGormRepository) AddProductToCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity float64, unitPrice float64) error {
	if err := validateCartItemAmount(quantity); err != nil {
		return err
	}
	if unitPrice <= 0 {
		return except.ErrCartInvalidPrice
	}

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart, err := r.getCartByUserIDWithDB(ctx, tx, userID)
		if err != nil {
			return err
		}

		updated, err := r.updateCartProductWithDB(ctx, tx, productID, quantity, unitPrice, cart)
		if err != nil {
			return err
		}
		if updated {
			return nil
		}

		lineTotal := quantity * unitPrice
		newProductItem := CartItem{
			ID:          uuid.New(),
			CartID:      cart.ID,
			ProductID:   productID,
			Quantity:    quantity,
			BasePrice:   unitPrice,
			TotalPrice:  lineTotal,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Components:  []byte("{}"),
		}
		if err := tx.WithContext(ctx).Create(&newProductItem).Error; err != nil {
			return fmt.Errorf("%w: %v", except.ErrCartAddProduct, err)
		}

		return r.recalculateTotalAmountWithDB(ctx, tx, cart.ID)
	})
}

func (r *cartGormRepository) updateCartProductWithDB(ctx context.Context, db *gorm.DB, productID uuid.UUID, quantityDelta float64, unitPrice float64, cart *Cart) (bool, error) {
	var existingItem *CartItem
	for i, item := range cart.Items {
		if item.ProductID == productID {
			existingItem = &cart.Items[i]
			break
		}
	}

	if existingItem != nil {
		newQuantity := existingItem.Quantity + quantityDelta
		if newQuantity < -qtyEpsilon {
			return false, except.ErrCartInvalidQuantity
		}
		if unitPrice <= 0 {
			return false, except.ErrCartInvalidPrice
		}
		if newQuantity > maxCartItemAmount+qtyEpsilon {
			return false, except.ErrCartInvalidQuantity
		}

		existingItem.Quantity = newQuantity
		existingItem.BasePrice = unitPrice
		existingItem.TotalPrice = newQuantity * unitPrice
		existingItem.UpdatedAt = time.Now()

		if newQuantity <= qtyEpsilon {
			if err := db.WithContext(ctx).
				Where("cart_id = ? AND product_id = ?", cart.ID, productID).
				Delete(&CartItem{}).Error; err != nil {
				return false, except.ErrCartUpdated
			}
		} else if err := db.WithContext(ctx).Save(existingItem).Error; err != nil {
			return false, except.ErrCartUpdated
		}

		if err := r.recalculateTotalAmountWithDB(ctx, db, cart.ID); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (r *cartGormRepository) getCartByUserIDWithDB(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*Cart, error) {
	var cart Cart

	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCart := NewCart(userID)
			if err := db.WithContext(ctx).Create(newCart).Error; err != nil {
				return nil, except.ErrCartCreate
			}
			return newCart, nil
		}

		return nil, except.ErrCartFind
	}

	return &cart, nil
}

func (r *cartGormRepository) RemoveProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart, err := r.getCartByUserIDWithDB(ctx, tx, userID)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).
			Where("cart_id = ? AND product_id = ?", cart.ID, productID).
			Delete(&CartItem{}).Error; err != nil {
			return except.ErrCartRemoveProduct
		}

		return r.recalculateTotalAmountWithDB(ctx, tx, cart.ID)
	})
}

func (r *cartGormRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart, err := r.getCartByUserIDWithDB(ctx, tx, userID)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).
			Where("cart_id = ?", cart.ID).
			Delete(&CartItem{}).Error; err != nil {
			return except.ErrCartClear
		}

		if err := tx.WithContext(ctx).
			Model(&Cart{}).
			Where("id = ?", cart.ID).
			Update("total_amount", 0).Error; err != nil {
			return except.ErrCartClear
		}

		return nil
	})
}

func (r *cartGormRepository) recalculateTotalAmountWithDB(ctx context.Context, db *gorm.DB, cartID uuid.UUID) error {
	var total float64
	if err := db.WithContext(ctx).
		Model(&CartItem{}).
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&total).Error; err != nil {
		return except.ErrCartRecalculate
	}

	if err := db.WithContext(ctx).
		Model(&Cart{}).
		Where("id = ?", cartID).
		Update("total_amount", total).Error; err != nil {
		return except.ErrCartRecalculate
	}

	return nil
}

func (r *cartGormRepository) RecalculateTotalAmount(ctx context.Context, cartID uuid.UUID) error {
	return r.recalculateTotalAmountWithDB(ctx, r.DB, cartID)
}

func (r *cartGormRepository) UpdateCartProduct(ctx context.Context, productID uuid.UUID, quantityDelta float64, unitPrice float64, cart *Cart) (bool, error) {
	return r.updateCartProductWithDB(ctx, r.DB, productID, quantityDelta, unitPrice, cart)
}

func (r *cartGormRepository) GetCartByUserId(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	return r.getCartByUserIDWithDB(ctx, r.DB, userID)
}
