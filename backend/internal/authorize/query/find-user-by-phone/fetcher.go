package find_user_by_phone

import (
	"context"

	"github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type Fetcher struct {
	repo entity.RegistrationRepository
}

func NewFetcher(repo entity.RegistrationRepository) *Fetcher {
	return &Fetcher{
		repo: repo,
	}
}

func (h *Fetcher) Fetch(ctx context.Context, phone string) (*entity.User, error) {
	user, err := h.repo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return user, nil
}
