package config

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

var em = eventManager.New()

func GetEventManager() *eventManager.EventManager {
	return em
}

// TODO: ВАЖНО СОХРАНЯТЬ ПОРЯДОК ПОСЛЕДОВАТЕЛЬНОСТИ ДОБАВЛЕНИЯ СЛУШАТЕЛЕЙ
func InitEventManager(em *eventManager.EventManager, listeners *action.Listeners) {

	em.AddEventListener("menu.refreshed", listeners.SaveMenu.Listen)
	em.AddEventListener("menu.refreshed", listeners.SaveModifier.Listen)
	em.AddEventListener("menu.refreshed", listeners.SaveProduct.Listen)

	em.AddEventListener("order.provider.created", listeners.OrderProviderCreated.Listen)
	em.AddEventListener("order.provider.created", listeners.QueueOrderProviderCreated.Listen)
	em.AddEventListener("order.status.changed", listeners.OrderStatusChanged.Listen)
	em.AddEventListener("order.status.changed", listeners.QueueOrderStatusChanged.Listen)
	em.Start()
}

type RefreshMenuEventPayload struct{}
