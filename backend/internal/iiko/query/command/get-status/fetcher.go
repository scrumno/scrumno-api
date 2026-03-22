package get_status

import (
	"context"
	"fmt"
	"strings"

	commandStatus "github.com/scrumno/scrumno-api/internal/iiko/entity/command-status"
)

type Fetcher struct {
	repo commandStatus.Repository
}

func NewFetcher(repo commandStatus.Repository) *Fetcher {
	return &Fetcher{
		repo: repo,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, query Query) (*commandStatus.StatusResponse, error) {
	organizationID := strings.TrimSpace(query.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для получения статуса команды iiko")
	}

	correlationID := strings.TrimSpace(query.CorrelationID)
	if correlationID == "" {
		return nil, fmt.Errorf("не передан correlationId для получения статуса команды iiko")
	}

	return f.repo.GetStatus(ctx, commandStatus.GetStatusRequest{
		OrganizationID: organizationID,
		CorrelationID:  correlationID,
	})
}
