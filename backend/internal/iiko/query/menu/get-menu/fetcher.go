package get_menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/menu"
)

type Fetcher struct {
	repo menu.Repository
}

func NewFetcher(repo menu.Repository) *Fetcher {
	return &Fetcher{
		repo: repo,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, query Query) (*menu.Menu, error) {
	organizationID := strings.TrimSpace(query.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для получения меню iiko")
	}

	return f.repo.GetList(ctx, menu.GetListParams{
		OrganizationID: organizationID,
		StartRevision:  query.StartRevision,
	})
}
