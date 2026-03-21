package check_tokens_revoke

import (
	"context"
	"errors"

	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
)

type Fetcher struct {
	tokenRepo authorizeTokens.TokensRepository
}

func NewFetcher(tokenRepo authorizeTokens.TokensRepository) *Fetcher {
	return &Fetcher{
		tokenRepo: tokenRepo,
	}
}

func (h *Fetcher) Fetch(ctx context.Context, accessToken string) (bool, error) {

	active, err := h.tokenRepo.IsSessionActive(ctx, accessToken)
	if err != nil {
		return false, err
	}
	if !active {
		return false, errors.New("Сессия завершена")
	}
	return true, nil
}
