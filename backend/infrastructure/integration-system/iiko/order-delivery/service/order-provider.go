package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoHTTP "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/model"
	helpers "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type OrderProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewOrderProvider(config *iikoConfig.Config) *OrderProvider {
	return &OrderProvider{
		http:   iikoHTTP.NewClient(config),
		config: config,
	}
}
func (p *OrderProvider) SetOrder(ctx context.Context, builderBody any) (any, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/deliveries/create", base)

	bodyJson, err := helpers.CreateBody[*BuilderSetCommand](builderBody)
	if err != nil {
		return nil, err
	}

	helpers.Logger([]byte(bodyJson))

	resp, err := helpers.SendRequest(p.http, url, bodyJson, p.config.AccessToken)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, nil
	}

	helpers.Logger([]byte(resp))

	var response model.OrderSetResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
