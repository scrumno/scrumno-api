package update_product_to_cart

import (
	"context"

	cartRepo "github.com/scrumno/scrumno-api/internal/cart/entity"
)

type Handler struct {
	cr cartRepo.CartRepository
}

func NewHandler(cartRepo cartRepo.CartRepository) *Handler {
	return &Handler{
		cr: cartRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	cart, err := h.cr.GetCartByUserId(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	var unitPrice float64
	var quantityDelta float64
	found := false

	for _, item := range cart.Items {
		if item.ProductID == cmd.ProductID {
			unitPrice = item.BasePrice
			quantityDelta = cmd.Quantity - item.Quantity
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	_, err = h.cr.UpdateCartProduct(
		ctx,
		cmd.ProductID,
		quantityDelta,
		unitPrice,
		cart,
	)

	return err
}
