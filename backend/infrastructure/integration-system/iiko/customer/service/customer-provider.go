package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	model "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/customer/model"
	iikoHTTP "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http"
	helpers "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type CustomerProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewCustomerProvider(config *iikoConfig.Config) *CustomerProvider {
	p := &CustomerProvider{
		http:   iikoHTTP.NewClient(config),
		config: config,
	}

	return p
}

func (p *CustomerProvider) GetCustomer(ctx context.Context, builderBody any) (any, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/loyalty/iiko/customer/info", base)

	bodyJson, err := helpers.CreateBody[*BuilderCommand](builderBody)
	if err != nil {
		return nil, err
	}

	resp, err := helpers.SendRequest(p.http, url, bodyJson, p.config.AccessToken)
	if err != nil {
		errText := err.Error()
		if strings.Contains(errText, "status=400") && strings.Contains(errText, "Validation_IncorrectPhone") {
			return nil, nil
		}
		return nil, err
	}

	if len(resp) == 0 {
		return nil, nil
	}

	var response model.CustomerResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *CustomerProvider) SetCustomer(ctx context.Context, builderBody any) (any, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/loyalty/iiko/customer/create_or_update", base)

	bodyJson, err := helpers.CreateBody[*BuilderSetCommand](builderBody)
	if err != nil {
		return nil, err
	}

	resp, err := helpers.SendRequest(p.http, url, bodyJson, p.config.AccessToken)
	if err != nil {
		return nil, err
	}
	slog.Info("iiko SetCustomer response", "raw", string(resp))
	if len(resp) == 0 {
		return nil, nil
	}

	var response model.CustomerSetResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
