package get_refresh

import (
	"context"
	"errors"

	tokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
)

type Fetcher struct {
	tokenRepo  tokens.TokensRepository
}

func NewFetcher(tokenRepo tokens.TokensRepository) *Fetcher {
	return &Fetcher{
		tokenRepo: tokenRepo,
	}
}

func (h *Fetcher) Fetch(ctx context.Context, accessToken string) (bool, error) {

	active, err := h.tokenRepo.IsSessionActive(ctx, accessToken);
	if err != nil {
		return false, err
	}
	if !active {
		return false, errors.New("Сессия завершена")
	}
	return true, nil
}