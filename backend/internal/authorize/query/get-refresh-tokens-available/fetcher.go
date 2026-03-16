package get_refresh_tokens_available

import (
	"context"
	"github.com/scrumno/scrumno-api/shared/jwt"
	tokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
)

type Fetcher struct {
	tokenRepo  tokens.TokensRepository
	jwtManager *jwt.Manager
}

func NewFetcher(tokenRepo tokens.TokensRepository, jwtManager *jwt.Manager) *Fetcher {
	return &Fetcher{
		tokenRepo: tokenRepo,
		jwtManager: jwtManager,
	}
}

func (h *Fetcher) Fetch(ctx context.Context, refreshToken string) (*jwt.Claims, error) {

	claims, err := h.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	storedToken, err := h.tokenRepo.FindTokenPairBySessionId(ctx, claims.SessionID)
	if err != nil || storedToken == nil {
		return nil, err
	}

	return claims, nil
}