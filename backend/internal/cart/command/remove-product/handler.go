package remmove_product

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

func (h *Handler) Handle(ctx context.Context, cmd Command) (bool, error) {
	cart, err := h.cr.GetCartByUserId(ctx, cmd.UserID)
	if err != nil {
		return false, err
	}
	if cart == nil {
		return false, nil
	}

	if err := h.cr.RemoveProduct(ctx, cmd.UserID, cmd.ProductID); err != nil {
		return false, err
	}

	return true, nil
}
