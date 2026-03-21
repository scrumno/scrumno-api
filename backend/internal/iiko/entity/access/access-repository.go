package access

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AccessRepository инкапсулирует HTTP-запросы к API iiko (аутентификация и далее — по мере роста домена).
type AccessRepository struct {
	baseURL string
	client  *http.Client
}

func NewAccessRepository(baseURL string, client *http.Client) *AccessRepository {
	if client == nil {
		client = http.DefaultClient
	}
	return &AccessRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

type accessTokenRequest struct {
	APILogin    string `json:"apiLogin"`
	APIPassword string `json:"apiPassword,omitempty"`
}

type accessTokenResponse struct {
	Token string `json:"token"`
}

func (r *AccessRepository) PostAccessToken(ctx context.Context, apiLogin, apiPassword string) (string, error) {
	u := r.baseURL + "/api/1/access_token"
	body, err := json.Marshal(accessTokenRequest{
		APILogin:    apiLogin,
		APIPassword: apiPassword,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("iiko access_token: status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out accessTokenResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("decode access_token response: %w", err)
	}
	if out.Token == "" {
		return "", fmt.Errorf("iiko access_token: empty token in response")
	}
	return out.Token, nil
}
