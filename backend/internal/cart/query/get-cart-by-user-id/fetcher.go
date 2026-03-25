package get_cart_by_user_id

import (
	"context"

	cartRepo "github.com/scrumno/scrumno-api/internal/cart/entity"
)

type Fetcher struct {
	cr cartRepo.CartRepository
}

func NewFetcher(cartRepo cartRepo.CartRepository) *Fetcher {
	return &Fetcher{
		cr: cartRepo,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, q Query) (*cartRepo.Cart, error) {
	cart, err := f.cr.GetCartByUserId(ctx, q.UserID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}
