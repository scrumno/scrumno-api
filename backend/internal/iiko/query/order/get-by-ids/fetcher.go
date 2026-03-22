package get_by_ids

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrdersByIDsGetter interface {
	GetByIDs(ctx context.Context, request order.GetByIDRequest) (*order.TableOrdersResponse, error)
}

type Fetcher struct {
	repo OrdersByIDsGetter
}

func NewFetcher(repo OrdersByIDsGetter) *Fetcher {
	return &Fetcher{repo: repo}
}

func (f *Fetcher) Fetch(ctx context.Context, query Query) (*order.TableOrdersResponse, error) {
	organizationIDs := sanitizeIDs(query.OrganizationIDs)
	if len(organizationIDs) == 0 {
		return nil, fmt.Errorf("не переданы organizationIds для получения заказов iiko")
	}

	orderIDs := sanitizeIDs(query.OrderIDs)
	posOrderIDs := sanitizeIDs(query.PosOrderIDs)

	if len(orderIDs) == 0 && len(posOrderIDs) == 0 {
		return nil, fmt.Errorf("не переданы orderIds или posOrderIds для получения заказов iiko")
	}
	if len(orderIDs) > 0 && len(posOrderIDs) > 0 {
		return nil, fmt.Errorf("нельзя передавать одновременно orderIds и posOrderIds")
	}

	return f.repo.GetByIDs(ctx, order.GetByIDRequest{
		OrganizationIDs: organizationIDs,
		OrderIDs:        orderIDs,
		PosOrderIDs:     posOrderIDs,
	})
}

func sanitizeIDs(ids []string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
