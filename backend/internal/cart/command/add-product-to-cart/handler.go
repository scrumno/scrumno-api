package add_product_to_cart

import (
	"context"

	cartRepo "github.com/scrumno/scrumno-api/internal/cart/entity"
)

type Handler struct {
	cr cartRepo.CartRepository
}

func NewHandler(cartRepository cartRepo.CartRepository) *Handler {
	return &Handler{
		cr: cartRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	if err := h.cr.AddProductToCart(
		ctx,
		cmd.UserID,
		cmd.ProductID,
		cmd.Quantity,
		cmd.BasePrice,
	); err != nil {
		return err
	}

	return nil
}
