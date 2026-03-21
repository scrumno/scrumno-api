package command_status

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
	GetStatus(ctx context.Context, request GetStatusRequest) (*StatusResponse, error)
}

type repository struct {
	baseURL string
	client  *http.Client
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

func (r *repository) GetStatus(ctx context.Context, request GetStatusRequest) (*StatusResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса статуса команды iiko: %w", err)
	}

	u := r.baseURL + "/api/1/commands/status"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос статуса команды iiko: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса статуса команды iiko: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ статуса команды iiko: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка статуса команды iiko: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out StatusResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ статуса команды iiko: %w", err)
	}

	return &out, nil
}
