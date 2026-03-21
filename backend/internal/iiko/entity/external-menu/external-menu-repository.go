package external_menu

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
	GetList(ctx context.Context, params GetListParams) (*MenusDataResponse, error)
	GetByID(ctx context.Context, params GetByIDParams) (*ByIDResponse, error)
}

type repository struct {
	baseURL string
	client  *http.Client
}

type GetListParams struct {
	OrganizationIDs []string
}

type getListRequest struct {
	OrganizationIDs []string `json:"organizationIds"`
}

type GetByIDParams struct {
	Request MenuRequest
}

func NewRepository(baseURL string, client *http.Client) Repository {
	if client == nil {
		client = http.DefaultClient
	}

	return &repository{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (r *repository) GetList(ctx context.Context, params GetListParams) (*MenusDataResponse, error) {
	organizationIDs := sanitizeOrganizationIDs(params.OrganizationIDs)
	if len(organizationIDs) == 0 {
		return nil, fmt.Errorf("не переданы organizationIds для запроса внешнего меню iiko")
	}

	body, err := json.Marshal(getListRequest{OrganizationIDs: organizationIDs})
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса внешнего меню iiko: %w", err)
	}

	raw, err := r.post(ctx, "/api/2/menu", body)
	if err != nil {
		return nil, err
	}

	var out MenusDataResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ внешнего меню iiko: %w", err)
	}

	return &out, nil
}

func (r *repository) GetByID(ctx context.Context, params GetByIDParams) (*ByIDResponse, error) {
	reqBody := params.Request
	reqBody.OrganizationIDs = sanitizeOrganizationIDs(reqBody.OrganizationIDs)
	if strings.TrimSpace(reqBody.ExternalMenuID) == "" {
		return nil, fmt.Errorf("не передан externalMenuId для запроса внешнего меню iiko")
	}
	if len(reqBody.OrganizationIDs) == 0 {
		return nil, fmt.Errorf("не переданы organizationIds для запроса внешнего меню iiko")
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса внешнего меню по id: %w", err)
	}

	raw, err := r.post(ctx, "/api/2/menu/by_id", body)
	if err != nil {
		return nil, err
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ внешнего меню по id: %w", err)
	}

	var formatVersion *int
	if v, ok := payload["formatVersion"].(float64); ok {
		asInt := int(v)
		formatVersion = &asInt
	}

	return &ByIDResponse{
		FormatVersion: formatVersion,
		Raw:           raw,
		Payload:       payload,
	}, nil
}

func (r *repository) post(ctx context.Context, endpoint string, body []byte) ([]byte, error) {
	u := r.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ %s: %w", endpoint, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка запроса %s: статус %d: %s", endpoint, resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	return raw, nil
}

func sanitizeOrganizationIDs(ids []string) []string {
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
