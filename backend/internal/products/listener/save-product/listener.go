package save_product

import (
	"context"
	"log/slog"

	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
	saveProductHandler "github.com/scrumno/scrumno-api/internal/products/command/save-product"
)

type Listener struct {
	handler *saveProductHandler.Handler
}

func NewListener(handler *saveProductHandler.Handler) *Listener {
	return &Listener{handler: handler}
}

func (l *Listener) Listen(payload any) {
	if payload == nil {
		slog.Info("save-product: <nil payload>")
		return
	}

	if _, ok := payload.(payloadMenuModel.RefreshMenuSuccessPayload); !ok {
		slog.Info("save-product: <invalid payload>")
		return
	}

	var cmd saveProductHandler.Command
	cmd.Products = payload.(payloadMenuModel.RefreshMenuSuccessPayload).Products

	err := l.handler.Handle(context.Background(), cmd)
	if err != nil {
		slog.Error("save-product: <error>", "error", err)
	}

}
