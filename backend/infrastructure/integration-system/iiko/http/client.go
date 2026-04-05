package httpclient

import (
	"net/http"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoMiddleware "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/http/middleware"
)

func NewClient(cfg *iikoConfig.Config) *http.Client {
	refresh := NewTokenRefresher(cfg)
	return &http.Client{
		Transport: iikoMiddleware.NewAuthRefreshRoundTripper(
			http.DefaultTransport,
			&cfg.AccessToken,
			refresh,
		),
	}
}
