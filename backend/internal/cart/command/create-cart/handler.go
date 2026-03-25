package create_cart

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

func (h *Handler) Handle(ctx context.Content, cmd Command) (*Cart, error) {
	newCart := h.cr.NewCart(cmd.UserID)

	if err := h.cr.Create(newCart).Error; err != nil {
		return nil, err
	}

	return cart, nil;
}