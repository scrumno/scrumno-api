package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
)

// NewTokenRefresher возвращает функцию, которая запрашивает новый access token у iiko.
func NewTokenRefresher(cfg *iikoConfig.Config) func(ctx context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		if strings.TrimSpace(cfg.Login) == "" {
			return "", fmt.Errorf("логин iiko пуст (IIKO_LOGIN)")
		}

		base := strings.TrimRight(cfg.BaseURL, "/")
		url := fmt.Sprintf("%s/api/1/access_token", base)

		reqBody, err := json.Marshal(map[string]string{
			"apiLogin": cfg.Login,
		})
		if err != nil {
			return "", err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		// Важно: не используем клиент с middleware, чтобы избежать рекурсии.
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

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return "", fmt.Errorf("обновление токена iiko не удалось: status=%d body=%s", resp.StatusCode, string(respBody))
		}

		var tokenResponse struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
			return "", err
		}
		if strings.TrimSpace(tokenResponse.Token) == "" {
			return "", fmt.Errorf("обновление токена iiko вернуло пустой токен")
		}

		return tokenResponse.Token, nil
	}
}

