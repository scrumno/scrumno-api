package authorize_service

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/access"
)

type Service struct {
	access *access.AccessRepository
}

func NewService(accessRepo *access.AccessRepository) *Service {
	return &Service{access: accessRepo}
}

func (s *Service) GetAccessToken(ctx context.Context, apiLogin, apiPassword string) (string, error) {
	token, err := s.access.PostAccessToken(ctx, apiLogin, apiPassword)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}
