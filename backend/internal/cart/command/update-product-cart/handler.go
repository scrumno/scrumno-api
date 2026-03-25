package update_product_to_cart

import (
	"github.com/google/uuid"
	"context"
	"errors"

	cartRepo "github.com/scrumno/scrumno-api/internal/cart/entity"
	except "github.com/scrumno/scrumno-api/shared/exception/cart"
)

type Handler struct {
	cr cartRepo.CartRepository
}

func NewHandler(cartRepo cartRepo.CartRepository) *Handler {
	return &Handler{
		cr: cartRepo,
	}
}

func (h *Handler) Handler(ctx context.Context, cmd Command) error {
	cart, err := h.cr.GetCartByUserId(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	var productPrice float64; 
	found := false
    
    for _, product := range cart.Items {
        if product.ID == cmd.ProductID {
            productPrice = float64(product.UnitPrice) * float64(cmd.Quantity)
            found = true
            break
        }
    }

	if !found {
        return nil
    }

	isUpdated, err := h.cr.UpdateCartProduct(
		ctx,
		cmd.ProductID,
		cmd.Quantity,
		productPrice,
		cart,
	)

	if err != nil {
		return err
	}

	if !isUpdated {
		return except.ErrCartUpdated
	}

	return nil
}