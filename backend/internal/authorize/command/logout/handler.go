package logout

import (
	"context"

	"github.com/google/uuid"
    tokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
)

type Handler struct {
	tokenRepo tokens.TokensRepository
}

func NewHandler(tokenRepo tokens.TokensRepository) *Handler {
	return &Handler{
		tokenRepo: tokenRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, userId uuid.UUID) error {
	return h.tokenRepo.RevokeTokensByUserSessionId(ctx, userId)
}