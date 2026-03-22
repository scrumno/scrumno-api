package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/scrumno/scrumno-api/infra/integration-system/iiko/config"
	"github.com/scrumno/scrumno-api/infra/integration-system/iiko/order/model"
	"github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"
)

type orderProvider struct {
	client *http.Client
	config *config.Config
}

func NewOrderProvider(client *http.Client, config *config.Config) interfaces.OrderProvider {
	return &orderProvider{
		client: client,
		config: config,
	}
}

func (p *orderProvider) Create(order any) (any, error) {
	_, ok := order.(model.CreateOrderRequest)
	if !ok {
		return nil, errors.New("order: expected CreateOrderRequest")
	}

	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Post(p.config.BaseURL+"/api/1/deliveries/create", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return CreateOrderResponse{
		StatusCode: res.StatusCode,
	}, nil
}

type CreateOrderResponse struct {
	StatusCode int    `json:"statusCode"`
	ID         string `json:"id,omitempty"`
}
