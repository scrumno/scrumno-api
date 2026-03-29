package clear_cart

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
	err := h.cr.ClearCart(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	return nil
}
