package get_by_pos_order_ids

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderByPosOrderIDsGetter interface {
	GetByPosOrderIDs(ctx context.Context, request order.GetByIDRequest) (*order.TableOrdersResponse, error)
}

type Fetcher struct {
	repo OrderByPosOrderIDsGetter
}

func NewFetcher(repo OrderByPosOrderIDsGetter) *Fetcher {
	return &Fetcher{
		repo: repo,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, query Query) (*order.TableOrdersResponse, error) {
	posOrderIDs := make([]string, 0, len(query.PosOrderIDs))
	for _, id := range query.PosOrderIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		posOrderIDs = append(posOrderIDs, trimmed)
	}
	if len(posOrderIDs) == 0 {
		return nil, fmt.Errorf("не переданы posOrderIds для получения заказов iiko")
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
		return nil, fmt.Errorf("не переданы organizationIds для получения заказов iiko")
	}

	return f.repo.GetByPosOrderIDs(ctx, order.GetByIDRequest{
		OrganizationIDs: organizationIDs,
		PosOrderIDs:     posOrderIDs,
	})
}
