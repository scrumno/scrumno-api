package orders

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	payOrderDraft "github.com/scrumno/scrumno-api/internal/orders/command/pay-order-draft"
)

type PayOrderDraftAction struct {
	Handler *payOrderDraft.Handler
}

type PayOrderDraftRequest struct {
	DraftID uuid.UUID `json:"draft_id"`
}

type PayOrderDraftResponse struct {
	IsSuccess     bool      `json:"is_success"`
	DraftID       uuid.UUID `json:"draft_id,omitempty"`
	OrderID       uuid.UUID `json:"order_id,omitempty"`
	CanSubscribe  bool      `json:"can_subscribe"`
	ProviderState string    `json:"provider_state,omitempty"`
	Error         string    `json:"error,omitempty"`
}

func NewPayOrderDraftAction(handler *payOrderDraft.Handler) *PayOrderDraftAction {
	return &PayOrderDraftAction{Handler: handler}
}

func (a *PayOrderDraftAction) Action(w http.ResponseWriter, r *http.Request) {
	var req PayOrderDraftRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, PayOrderDraftResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if req.DraftID == uuid.Nil {
		utils.JSONResponse(w, PayOrderDraftResponse{
			IsSuccess: false,
			Error:     "Поле draft_id обязательно",
		}, http.StatusBadRequest)
		return
	}

	result := a.Handler.Handle(r.Context(), payOrderDraft.Command{
		DraftID: req.DraftID,
	})
	status := http.StatusOK
	if !result.IsSuccess {
		status = http.StatusBadRequest
	}

	utils.JSONResponse(w, PayOrderDraftResponse{
		IsSuccess:     result.IsSuccess,
		DraftID:       result.DraftID,
		OrderID:       result.OrderID,
		CanSubscribe:  result.CanSubscribe,
		ProviderState: result.ProviderState,
		Error:         result.Error,
	}, status)
}

func (a *PayOrderDraftAction) GetInputType() reflect.Type {
	return reflect.TypeOf(PayOrderDraftRequest{})
}
