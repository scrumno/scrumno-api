package create_authorize_tokens

import (
	"context"
	"time"
	"github.com/google/uuid"

    tokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	jwt "github.com/scrumno/scrumno-api/shared/jwt"
)

type Handler struct {
	tokenRepo tokens.TokensRepository
	jwtManager *jwt.Manager
}

func NewHandler(
	tokenRepo tokens.TokensRepository, 
	jwtManager *jwt.Manager,
) *Handler {
	return &Handler{
		tokenRepo: tokenRepo, 
		jwtManager: jwtManager,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (string, string, int64, error) {

	if err := h.tokenRepo.RevokeTokensByUserSessionId(ctx, cmd.UserID); err != nil {
		return "", "", 0, err
	}

	sessionID := uuid.New()

	pair, err := h.jwtManager.GenerateTokenPair(
		cmd.UserID.String(), 
		cmd.Phone, 
		sessionID.String(),
	)
	if err != nil {
		return "", "", 0, err
	}

	jwtToken := tokens.NewAuthorizeToken(
		sessionID, 
		cmd.UserID, 
		pair.RefreshToken, 
		h.jwtManager.RefreshExpiresAtUnix(), 
		time.Now(),
	)

	err = h.tokenRepo.Create(ctx, jwtToken)
	if err != nil {
		return "", "", 0, err
	}

	return pair.RefreshToken, pair.AccessToken, pair.ExpiresIn, nil
}
