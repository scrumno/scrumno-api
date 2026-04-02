package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	model "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/customer/model"
	iikoMiddleware "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http/middleware"
	helpers "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type CustomerProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewCustomerProvider(config *iikoConfig.Config) *CustomerProvider {
	p := &CustomerProvider{
		http:   &http.Client{},
		config: config,
	}

	// Включаем 401 -> refresh -> retry для запросов, сделанных этим провайдером.
	p.http = &http.Client{
		Transport: iikoMiddleware.NewAuthRefreshRoundTripper(
			http.DefaultTransport,
			&p.config.AccessToken,
			p.refreshAccessToken,
		),
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

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.AccessToken))

	resp, err := p.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"запрос клиента iiko не удался: status=%d body=%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var response model.CustomerResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *CustomerProvider) SetCustomer(ctx context.Context, builderBody any) (any, error) {
	body, ok := builderBody.(*BuilderSetCommand)
	if !ok {
		return nil, fmt.Errorf("Невалидный формат body, должен быть BuilderSetCommand, %T", builderBody)
	}

	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/loyalty/iiko/customer/create_or_update", base)

	bodyJson, err := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}

	// fmt.Println(req)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.AccessToken))

	resp, err := p.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"запрос клиента iiko не удался: status=%d body=%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var response model.CustomerSetResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *CustomerProvider) refreshAccessToken(ctx context.Context) (string, error) {
	if strings.TrimSpace(p.config.Login) == "" {
		return "", fmt.Errorf("логин iiko пуст (IIKO_LOGIN)")
	}

	// Token endpoint format:
	// POST <base>/api/1/access_token
	// body: {"apiLogin":"<login>"}
	// response: {"token":"<access_token>"}
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/access_token", base)

	reqBody, err := json.Marshal(map[string]string{
		"apiLogin": p.config.Login,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// ВАЖНО: не используем p.http здесь, иначе мы попадем в ту же middleware.
	client := &http.Client{Transport: http.DefaultTransport}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("обновление токена iiko не удалось: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var tokenResponse struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
		return "", err
	}
	if tokenResponse.Token == "" {
		return "", fmt.Errorf("обновление токена iiko вернуло пустой токен")
	}

	return tokenResponse.Token, nil
}
