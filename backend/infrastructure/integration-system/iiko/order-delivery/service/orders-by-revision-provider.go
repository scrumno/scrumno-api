package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoHTTP "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http"
	helpers "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type OrdersByRevisionRequest struct {
	StartRevision   int64    `json:"startRevision"`
	OrganizationIDs []string `json:"organizationIds"`
	SourceKeys      []string `json:"sourceKeys,omitempty"`
}

type OrdersByRevisionResponse struct {
	MaxRevision           int64 `json:"maxRevision"`
	OrdersByOrganizations []struct {
		OrganizationID string `json:"organizationId"`
		Orders         []struct {
			ID    string `json:"id"`
			Order struct {
				Status         string `json:"status"`
				IsDeleted      *bool  `json:"isDeleted"`
				CompleteBefore string `json:"completeBefore"`
				Items          []struct {
					Amount float64 `json:"amount"`
				} `json:"items"`
			} `json:"order"`
		} `json:"orders"`
	} `json:"ordersByOrganizations"`
}

type OrdersByRevisionProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewOrdersByRevisionProvider(config *iikoConfig.Config) *OrdersByRevisionProvider {
	return &OrdersByRevisionProvider{
		http:   iikoHTTP.NewClient(config),
		config: config,
	}
}

func (p *OrdersByRevisionProvider) GetOrdersByRevision(ctx context.Context, startRevision int64, sourceKeys []string) (*OrdersByRevisionResponse, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/deliveries/by_revision", base)

	body := OrdersByRevisionRequest{
		StartRevision:   startRevision,
		OrganizationIDs: []string{p.config.OrganizationID.String()},
		SourceKeys:      sourceKeys,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := helpers.SendRequest(p.http, url, payload, p.config.AccessToken)
	if err != nil {
		return nil, err
	}

	var parsed OrdersByRevisionResponse
	if err := json.Unmarshal(resp, &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}
