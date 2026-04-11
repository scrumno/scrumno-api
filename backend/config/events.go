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
	if listeners.SaveMenu != nil {
		em.AddEventListener("menu.refreshed", listeners.SaveMenu.Listen)
	}
	if listeners.SaveModifier != nil {
		em.AddEventListener("menu.refreshed", listeners.SaveModifier.Listen)
	}
	if listeners.SaveProduct != nil {
		em.AddEventListener("menu.refreshed", listeners.SaveProduct.Listen)
	}
	if listeners.OrderProviderCreated != nil {
		em.AddEventListener("order.provider.created", listeners.OrderProviderCreated.Listen)
	}
	if listeners.QueueOrderProviderCreated != nil {
		em.AddEventListener("order.provider.created", listeners.QueueOrderProviderCreated.Listen)
	}
	if listeners.OrderStatusChanged != nil {
		em.AddEventListener("order.status.changed", listeners.OrderStatusChanged.Listen)
	}
	if listeners.QueueOrderStatusChanged != nil {
		em.AddEventListener("order.status.changed", listeners.QueueOrderStatusChanged.Listen)
	}
	em.Start()
}

type RefreshMenuEventPayload struct{}
