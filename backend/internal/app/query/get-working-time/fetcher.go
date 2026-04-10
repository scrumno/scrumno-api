package get_working_time

import (
	"context"

	appConfig "github.com/scrumno/scrumno-api/internal/app/entity/app-config"
)

type WorkingHours struct {
	OpenAt  string `json:"open_at"`
	CloseAt string `json:"close_at"`
}

type Fetcher struct {
	appConfigRepository appConfig.AppConfigRepository
}

func NewFetcher(appConfigRepository appConfig.AppConfigRepository) *Fetcher {
	return &Fetcher{appConfigRepository: appConfigRepository}
}

func (h *Fetcher) Fetch(ctx context.Context) (WorkingHours, error) {
	openAt, closeAt, err := h.appConfigRepository.GetWorkingHours(ctx)
	if err != nil {
		return WorkingHours{}, err
	}
	return WorkingHours{
		OpenAt:  openAt,
		CloseAt: closeAt,
	}, nil
}
