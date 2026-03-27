package refresh_menu

import (
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

type Handler struct {
	provider     interfaces.MenuProvider
	eventManager *eventManager.EventManager
}

func NewHandler(provider interfaces.MenuProvider, eventManager *eventManager.EventManager) *Handler {
	return &Handler{
		provider:     provider,
		eventManager: eventManager,
	}
}

func (h *Handler) Handle() any {
	menu, err := h.provider.GetMenu()
	if err != nil {
		return err
	}

	if h.eventManager != nil {
		h.eventManager.EmitEvent("menu.refreshed", menu)
	}
	return menu
}
