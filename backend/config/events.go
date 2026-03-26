package config

import (
	refreshMenuCmd "github.com/scrumno/scrumno-api/internal/menu/command/refresh-menu"
	refreshMenuListener "github.com/scrumno/scrumno-api/internal/menu/listener/refresh-menu"
	menuInterfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

var em = eventManager.New()

func GetEventManager() *eventManager.EventManager {
	return em
}

func InitEventManager(em *eventManager.EventManager, menuProvider menuInterfaces.MenuProvider) {
	refreshMenuHandler := refreshMenuCmd.NewHandler(menuProvider, em)
	refreshMenuEventListener := refreshMenuListener.NewListener(refreshMenuHandler)

	em.AddEventListener("menu.refreshed", refreshMenuEventListener.Listen)
	em.Start()
}

type RefreshMenuEventPayload struct{}
