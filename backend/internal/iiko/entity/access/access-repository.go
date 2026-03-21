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

func (r *AccessRepository) PostAccessToken(ctx context.Context, apiLogin, apiPassword string) (*AccessToken, error) {
	if strings.TrimSpace(apiLogin) == "" {
		return nil, fmt.Errorf("логин iiko (apiLogin) не задан")
	}

	u := r.baseURL + "/api/1/access_token"
	body, err := json.Marshal(accessTokenRequest{
		APILogin:    apiLogin,
		APIPassword: apiPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать тело запроса access_token: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос access_token: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса access_token: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ access_token: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("iiko access_token: статус %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var out accessTokenResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ access_token: %w", err)
	}
	if out.Token == "" {
		return nil, fmt.Errorf("iiko access_token: пустой token в ответе")
	}

	return &AccessToken{Token: out.Token}, nil
}
