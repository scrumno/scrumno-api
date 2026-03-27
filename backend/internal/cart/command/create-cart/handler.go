package create_cart

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

func (h *Handler) Handle(ctx context.Context, cmd Command) (*cartRepo.Cart, error) {
	newCart := cartRepo.NewCart(cmd.UserID)

	if err := h.cr.Create(ctx, newCart); err != nil {
		return nil, err
	}

	return newCart, nil
}
