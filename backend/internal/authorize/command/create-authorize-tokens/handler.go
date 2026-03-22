package create_authorize_tokens

import (
	"context"
	"time"

	"github.com/google/uuid"

	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	jwt "github.com/scrumno/scrumno-api/shared/jwt"
)

type Handler struct {
	tokenRepo  authorizeTokens.TokensRepository
	jwtManager *jwt.Manager
}

func NewHandler(
	tokenRepo authorizeTokens.TokensRepository,
	jwtManager *jwt.Manager,
) *Handler {
	return &Handler{
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (string, string, int64, error) {

	if cmd.RevokePreviousToken && cmd.SessionID != "" {
		err := h.tokenRepo.RevokeTokenBySessionId(ctx, cmd.SessionID)
		if err != nil {
			return "", "", 0, err
		}
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

	jwtToken := authorizeTokens.NewAuthorizeToken(
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

	return pair.AccessToken, pair.RefreshToken, pair.ExpiresIn, nil
}
