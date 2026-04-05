package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoHTTP "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http"
	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
)

type MenuProvider struct {
	http   *http.Client
	config *iikoConfig.Config
}

func NewMenuProvider(config *iikoConfig.Config) *MenuProvider {
	p := &MenuProvider{
		http:   iikoHTTP.NewClient(config),
		config: config,
	}

	return p
}

func (p *MenuProvider) GetMenu() (any, error) {
	base := strings.TrimRight(p.config.BaseURL, "/")
	url := fmt.Sprintf("%s/api/1/nomenclature", base)

	requestBody, err := json.Marshal(map[string]string{
		"organizationId":  p.config.OrganizationID.String(),
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
