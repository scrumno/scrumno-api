package config

import (
	refreshMenuCmd "github.com/scrumno/scrumno-api/internal/menu/command/refresh-menu"
	refreshMenuListener "github.com/scrumno/scrumno-api/internal/menu/listener/refresh-menu"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

var em = eventManager.New()

func GetEventManager() *eventManager.EventManager {
	return em
}

func InitEventManager(em *eventManager.EventManager) {

	RefreshMenuHandler := refreshMenuCmd.NewHandler()
	RefreshMenuListener := refreshMenuListener.NewListener(RefreshMenuHandler)

	em.AddEventListener("menu.refresh", RefreshMenuListener.Listen)
	em.Start()
}

type RefreshMenuEventPayload struct {
}
