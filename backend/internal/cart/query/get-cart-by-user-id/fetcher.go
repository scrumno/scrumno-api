package get_cart_by_user_id

import (
	"github.com/google/uuid"
	"context"
	"errors"

	cartRepo "github.com/scrumno/scrumno-api/internal/cart/entity"
	except "github.com/scrumno/scrumno-api/shared/exception/cart"
)

type Fetcher struct {
	cr cartRepo.CartRepository
}

func NewFetcher(cartRepo cartRepo.CartRepository) *Fetcher {
	return &Fetcher{
		cr: cartRepo,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, q Query) (*Cart, error) {
	cart, err := h.cr.GetCartByUserId(ctx, q.UserID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}