package menu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Repository interface {
	GetList(ctx context.Context, params GetListParams) (*Menu, error)
}

type menuRepository struct {
	baseURL string
	client  *http.Client
}

type getListRequest struct {
	OrganizationID string `json:"organizationId"`
	StartRevision  *int64 `json:"startRevision,omitempty"`
}

type GetListParams struct {
	OrganizationID string
	StartRevision  *int64
}

func NewRepository(baseURL string, client *http.Client) Repository {
	if client == nil {
		client = http.DefaultClient
	}
	return &menuRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (r *menuRepository) GetList(ctx context.Context, params GetListParams) (*Menu, error) {
	organizationID := strings.TrimSpace(params.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для запроса меню iiko")
	}

	u := r.baseURL + "/api/1/nomenclature"
	body, err := json.Marshal(getListRequest{
		OrganizationID: organizationID,
		StartRevision:  params.StartRevision,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса меню iiko: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос меню iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса меню iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ меню iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка запроса меню iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out Menu
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ меню iiko: %w", err)
	}

	return &out, nil
}
