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
	iikoMiddleware "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http/middleware"
	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
)

type MenuProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewMenuProvider(config *iikoConfig.Config) *MenuProvider {
	p := &MenuProvider{
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

func (p *MenuProvider) refreshAccessToken(ctx context.Context) (string, error) {
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

func (p *MenuProvider) GetMenu() (any, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/nomenclature", base)

	requestBody, err := json.Marshal(map[string]string{
		"organizationId":  p.config.OrganizationID,
		"terminalGroupId": p.config.TerminalGroupID,
	})
	if err != nil {
		return payloadMenuModel.RefreshMenuSuccessPayload{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return payloadMenuModel.RefreshMenuSuccessPayload{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.AccessToken))

	resp, err := p.http.Do(req)
	if err != nil {
		return payloadMenuModel.RefreshMenuSuccessPayload{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return payloadMenuModel.RefreshMenuSuccessPayload{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return payloadMenuModel.RefreshMenuSuccessPayload{}, fmt.Errorf(
			"запрос меню iiko не удался: status=%d body=%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var response payloadMenuModel.RefreshMenuSuccessPayload

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return payloadMenuModel.RefreshMenuSuccessPayload{}, err
	}

	return response, nil
}
