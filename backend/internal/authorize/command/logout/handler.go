package logout

import (
	"context"

	"github.com/google/uuid"
	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
)

type Handler struct {
	tokenRepo authorizeTokens.TokensRepository
}

func NewHandler(tokenRepo authorizeTokens.TokensRepository) *Handler {
	return &Handler{
		tokenRepo: tokenRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, userId uuid.UUID) error {
	return h.tokenRepo.RevokeTokensByUserSessionId(ctx, userId)
}
