package iiko

import (
	"context"
	"net/http"
	"time"

	setAccess "github.com/scrumno/scrumno-api/internal/iiko/command/auth/set-access"
	iikoconfig "github.com/scrumno/scrumno-api/internal/iiko/config"
	"github.com/scrumno/scrumno-api/internal/iiko/entity/access"
	iikoMiddleware "github.com/scrumno/scrumno-api/internal/iiko/middleware"
	authorizeService "github.com/scrumno/scrumno-api/internal/iiko/service/authorize-service"
	tokenProvider "github.com/scrumno/scrumno-api/internal/iiko/service/token-provider"
)

type Container struct {
	SetAccess *setAccess.Handler
}

func NewContainer(cfg *iikoconfig.Config) *Container {
	if cfg == nil {
		c := iikoconfig.Load()
		cfg = &c
	}

	var accessRepo *access.AccessRepository
	provider := tokenProvider.NewProvider(func(ctx context.Context) (string, error) {
		return accessRepo.PostAccessToken(ctx, cfg.Login, cfg.Password)
	})

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: iikoMiddleware.NewAuthRetryTransport(
			http.DefaultTransport,
			provider,
		),
	}

	// repository
	accessRepo = access.NewAccessRepository(cfg.BaseURL, httpClient)

	// service
	authSvc := authorizeService.NewService(accessRepo)

	// command
	setAccessHandler := setAccess.NewHandler(authSvc, cfg.Login, cfg.Password)

	return &Container{
		SetAccess: setAccessHandler,
	}
}
