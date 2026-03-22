package set_access

import (
	"context"

	authorizeService "github.com/scrumno/scrumno-api/internal/iiko/service/authorize-service"
)

type Handler struct {
	svc      *authorizeService.Service
	apiLogin string
	apiPass  string
}

func NewHandler(svc *authorizeService.Service, apiLogin, apiPassword string) *Handler {
	return &Handler{
		svc:      svc,
		apiLogin: apiLogin,
		apiPass:  apiPassword,
	}
}

func (h *Handler) Handle(ctx context.Context) (token string, err error) {
	return h.svc.GetAccessToken(ctx, h.apiLogin, h.apiPass)
}
