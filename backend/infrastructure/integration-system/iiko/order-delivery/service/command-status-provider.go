package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoHTTP "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http"
	helpers "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type CommandStatusProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

type CommandStatusRequest struct {
	OrganizationID string    `json:"organizationId"`
	CorrelationID  uuid.UUID `json:"correlationId"`
}

type CommandStatusResponse struct {
	State       string  `json:"state"`
	ErrorReason *string `json:"errorReason,omitempty"`
}

func NewCommandStatusProvider(config *iikoConfig.Config) *CommandStatusProvider {
	return &CommandStatusProvider{
		http:   iikoHTTP.NewClient(config),
		config: config,
	}
}

func (p *CommandStatusProvider) GetStatus(ctx context.Context, correlationID uuid.UUID) (*CommandStatusResponse, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/commands/status", base)

	payload, err := json.Marshal(CommandStatusRequest{
		OrganizationID: p.config.OrganizationID.String(),
		CorrelationID:  correlationID,
	})
	if err != nil {
		return nil, err
	}

	resp, err := helpers.SendRequest(p.http, url, payload, p.config.AccessToken)
	if err != nil {
		return nil, err
	}

	var parsed CommandStatusResponse
	if err := json.Unmarshal(resp, &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}
