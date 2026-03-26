package refresh_menu

import (
	"fmt"
	"log/slog"

	refreshMenu "github.com/scrumno/scrumno-api/internal/menu/command/refresh-menu"
	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
)

type Listener struct {
	handler *refreshMenu.Handler
}

func NewListener(handler *refreshMenu.Handler) *Listener {
	return &Listener{handler: handler}
}

func (l *Listener) Listen(payload any) {
	if payload == nil {
		slog.Info("menu.refreshed: <nil payload>")
		return
	}

	if p, ok := payload.(payloadMenuModel.RefreshMenuSuccessPayload); ok {
		slog.Info(
			"menu.refreshed",
			"correlationId",
			p.CorrelationID,
			"groups",
			len(p.Groups),
			"productCategories",
			len(p.ProductCategories),
			"products",
			len(p.Products),
		)
		return
	}

	slog.Info("menu.refreshed", "payloadType", fmt.Sprintf("%T", payload))
}
