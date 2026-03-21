package get_by_id

import (
	"context"
	"fmt"
	"strings"

	externalMenu "github.com/scrumno/scrumno-api/internal/iiko/entity/external-menu"
)

type Fetcher struct {
	repo externalMenu.Repository
}

func NewFetcher(repo externalMenu.Repository) *Fetcher {
	return &Fetcher{
		repo: repo,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, query Query) (*externalMenu.ByIDResponse, error) {
	externalMenuID := strings.TrimSpace(query.ExternalMenuID)
	if externalMenuID == "" {
		return nil, fmt.Errorf("не передан externalMenuId для получения внешнего меню iiko")
	}

	organizationIDs := make([]string, 0, len(query.OrganizationIDs))
	for _, id := range query.OrganizationIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		organizationIDs = append(organizationIDs, trimmed)
	}
	if len(organizationIDs) == 0 {
		return nil, fmt.Errorf("не переданы organizationIds для получения внешнего меню iiko")
	}

	return f.repo.GetByID(ctx, externalMenu.GetByIDParams{
		Request: externalMenu.MenuRequest{
			ExternalMenuID:  externalMenuID,
			OrganizationIDs: organizationIDs,
			PriceCategoryID: query.PriceCategoryID,
			Version:         query.Version,
			Language:        query.Language,
			StartRevision:   query.StartRevision,
		},
	})
}
