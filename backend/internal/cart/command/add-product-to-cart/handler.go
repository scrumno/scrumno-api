package add_product_to_cart

import (
	"context"
	"errors"

	cartRepo "github.com/scrumno/scrumno-api/internal/cart/entity"
)

type Handler struct {
	cr cartRepo.CartRepository
}

func NewHandler(cartRepository cartRepo.cartRepository) *Handler {
	return &Handler{
		cr: cartRepository
	}
}

func (h *Handler) Handle(ctx context.Content, cmd Command) error {
	if err := h.cr.AddProductToCart(
		ctx, 
		cmd.UserID, 
		cmd.ProductID, 
		cmd.Qunatity, 
		cmd.BasePrice
	); err != nil {
		return err
	}

	return nil
}