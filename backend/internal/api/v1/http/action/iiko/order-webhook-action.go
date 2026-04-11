package iiko

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	processProviderWebhook "github.com/scrumno/scrumno-api/internal/orders/command/process-provider-webhook"
)

type OrderWebhookAction struct {
	Handler *processProviderWebhook.Handler
}

func NewOrderWebhookAction(handler *processProviderWebhook.Handler) *OrderWebhookAction {
	return &OrderWebhookAction{Handler: handler}
}

type webhookItem struct {
	EventType     string     `json:"eventType"`
	CorrelationID *uuid.UUID `json:"correlationId"`
	EventInfo     struct {
		ID             *uuid.UUID `json:"id"`
		CreationStatus string     `json:"creationStatus"`
		ErrorInfo      struct {
			Message string `json:"message"`
		} `json:"errorInfo"`
		Order *struct {
			Status string `json:"status"`
		} `json:"order"`
	} `json:"eventInfo"`
}

func (a *OrderWebhookAction) Action(w http.ResponseWriter, r *http.Request) {
	var payload []webhookItem
	if err := utils.DecodeJSONBody(r, &payload); err != nil {
		utils.JSONResponse(w, map[string]any{
			"is_success": false,
			"error":      err.Error(),
		}, http.StatusBadRequest)
		return
	}

	for _, item := range payload {
		var providerErr *string
		if item.EventInfo.ErrorInfo.Message != "" {
			providerErr = &item.EventInfo.ErrorInfo.Message
		}
		cmd := processProviderWebhook.Command{
			EventType:       item.EventType,
			CorrelationID:   item.CorrelationID,
			ProviderOrderID: item.EventInfo.ID,
			CreationStatus:  item.EventInfo.CreationStatus,
			Error:           providerErr,
		}
		if item.EventInfo.Order != nil {
			cmd.Status = item.EventInfo.Order.Status
		}
		_ = a.Handler.Handle(r.Context(), cmd)
	}

	utils.JSONResponse(w, map[string]any{
		"is_success": true,
	}, http.StatusOK)
}

func (a *OrderWebhookAction) GetInputType() reflect.Type {
	return reflect.TypeOf(webhookItem{})
}
