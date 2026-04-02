package save_menu

import (
	"context"
	"log/slog"

	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
	saveMenu "github.com/scrumno/scrumno-api/internal/menu/command/save-menu"
	categoryEntity "github.com/scrumno/scrumno-api/internal/menu/entity/category"
	sectionEntity "github.com/scrumno/scrumno-api/internal/menu/entity/section"
)

type Listener struct {
	handler *saveMenu.Handler
}

func NewListener(handler *saveMenu.Handler) *Listener {
	return &Listener{
		handler: handler,
	}
}

func (l *Listener) Listen(payload any) {
	if payload == nil {
		slog.Info("menu.saved: <nil payload>")
		return
	}

	menuPayload, ok := payload.(payloadMenuModel.RefreshMenuSuccessPayload)
	if !ok {
		slog.Info("menu.saved: <invalid payload>")
		return
	}

	var cmd saveMenu.Command
	for _, group := range menuPayload.Groups {
		cmd.Sections = append(cmd.Sections, sectionEntity.Section{
			ID:   group.ID,
			Name: group.Name,
		})
	}
	for _, category := range menuPayload.ProductCategories {
		cmd.Categories = append(cmd.Categories, categoryEntity.Category{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	if err := l.handler.Handle(context.Background(), cmd); err != nil {
		slog.Error("menu.saved: <error>", "error", err)
	}
}
