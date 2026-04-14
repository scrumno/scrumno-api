package refresh_menu

import (
	"log/slog"

	payloadMenuModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/model"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

type Handler struct {
	provider        interfaces.MenuProvider
	eventManager    *eventManager.EventManager
	snapshotService interfaces.SnapshotService
}

func NewHandler(provider interfaces.MenuProvider, eventManager *eventManager.EventManager, snapshotService interfaces.SnapshotService) *Handler {
	return &Handler{
		provider:        provider,
		eventManager:    eventManager,
		snapshotService: snapshotService,
	}
}

func (h *Handler) Handle() any {
	menu, err := h.provider.GetMenu()
	if err != nil {
		return err
	}

	const cmdName = "refresh-menu"

	menuForHash := menu
	if v, ok := menu.(payloadMenuModel.RefreshMenuSuccessPayload); ok {
		v.CorrelationID = ""
		menuForHash = v
	}

	isChanged, err := h.snapshotService.CheckAndSaveWithUploads(cmdName, menuForHash)
	if err != nil {
		// Если не удалось сравнить/сохранить снапшот, чтобы не "пропустить" обновления,
		// считаем, что меню могло измениться.
		slog.Error("refresh-menu: failed to save snapshot/photos", "error", err)
	} else {
		slog.Info("refresh-menu: menu changed", "changed", isChanged)
	}

	if h.eventManager != nil {
		h.eventManager.EmitEvent("menu.refreshed", menu)
	}

	return menu
}
