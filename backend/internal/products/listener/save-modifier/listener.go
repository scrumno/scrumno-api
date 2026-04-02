package save_modifier

import (
	"context"
	"log/slog"

	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
	saveModifier "github.com/scrumno/scrumno-api/internal/products/command/save-modifier"
)

type Listener struct {
	handler *saveModifier.Handler
}

func NewListener(handler *saveModifier.Handler) *Listener {
	return &Listener{handler: handler}
}

func (l *Listener) Listen(payload any) {
	if payload == nil {
		slog.Info("save-modifier: <nil payload>")
		return
	}

	menuPayload, ok := payload.(payloadMenuModel.RefreshMenuSuccessPayload)
	if !ok {
		slog.Info("save-modifier: <invalid payload>")
		return
	}

	var cmd saveModifier.Command
	for _, product := range menuPayload.Products {
		cmd.Modifiers = append(cmd.Modifiers, product.Modifiers...)
		cmd.Groups = append(cmd.Groups, product.GroupModifiers...)
		for _, group := range product.GroupModifiers {
			cmd.ChildModifiers = append(cmd.ChildModifiers, group.ChildModifiers...)
		}
	}

	if err := l.handler.Handle(context.Background(), cmd); err != nil {
		slog.Error("save-modifier: <error>", "error", err)
	}
}
