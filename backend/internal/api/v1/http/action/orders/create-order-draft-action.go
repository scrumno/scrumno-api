package orders

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
	userQuery "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	cartEntity "github.com/scrumno/scrumno-api/internal/cart/entity"
	cartQuery "github.com/scrumno/scrumno-api/internal/cart/query/get-cart-by-user-id"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
	createOrderDraft "github.com/scrumno/scrumno-api/internal/orders/command/create-order-draft"
)

type CreateOrderDraftAction struct {
	Handler     *createOrderDraft.Handler
	UserFetcher *userQuery.Fetcher
	CartFetcher *cartQuery.Fetcher
}

type CreateOrderDraftRequest struct {
	VenueID       uuid.UUID `json:"venue_id"`
	CustomerPhone string    `json:"phone"`
	Comment       string    `json:"comment"`
}

type CreateOrderDraftResponse struct {
	IsSuccess bool      `json:"is_success"`
	DraftID   uuid.UUID `json:"draft_id,omitempty"`
	Amount    float64   `json:"amount,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func NewCreateOrderDraftAction(handler *createOrderDraft.Handler, userFetcher *userQuery.Fetcher, cartFetcher *cartQuery.Fetcher) *CreateOrderDraftAction {
	return &CreateOrderDraftAction{
		Handler:     handler,
		UserFetcher: userFetcher,
		CartFetcher: cartFetcher,
	}
}

func (a *CreateOrderDraftAction) Action(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderDraftRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, CreateOrderDraftResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if req.VenueID == uuid.Nil {
		utils.JSONResponse(w, CreateOrderDraftResponse{
			IsSuccess: false,
			Error:     "Поле venue_id обязательно",
		}, http.StatusBadRequest)
		return
	}

	claims := middleware.ClaimsFromRequest(r)
	phone := strings.TrimSpace(req.CustomerPhone)
	if phone == "" && claims != nil {
		phone = claims.Phone
	}
	if phone == "" {
		utils.JSONResponse(w, CreateOrderDraftResponse{
			IsSuccess: false,
			Error:     "Укажите номер телефона",
		}, http.StatusBadRequest)
		return
	}

	user, err := a.UserFetcher.Fetch(r.Context(), phone)
	if err != nil || user == nil {
		utils.JSONResponse(w, CreateOrderDraftResponse{
			IsSuccess: false,
			Error:     "Пользователь не найден",
		}, http.StatusBadRequest)
		return
	}

	cart, err := a.CartFetcher.Fetch(r.Context(), cartQuery.Query{
		UserID: user.ID,
	})
	if err != nil || cart == nil || len(cart.Items) == 0 {
		utils.JSONResponse(w, CreateOrderDraftResponse{
			IsSuccess: false,
			Error:     "Корзина пользователя пуста",
		}, http.StatusBadRequest)
		return
	}

	result := a.Handler.Handle(r.Context(), createOrderDraft.Command{
		UserID:        user.ID,
		VenueID:       req.VenueID,
		CustomerPhone: phone,
		CustomerName:  draftCustomerDisplayName(user.FullName, phone),
		Comment:       optionalTrimmedString(req.Comment),
		SourceKey:     "app",
		CartItems:     mapCartItemsToDraftLines(cart),
	})
	if !result.IsSuccess {
		utils.JSONResponse(w, CreateOrderDraftResponse{
			IsSuccess: false,
			Error:     result.Error,
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, CreateOrderDraftResponse{
		IsSuccess: true,
		DraftID:   result.DraftID,
		Amount:    result.Amount,
	}, http.StatusOK)
}

func (a *CreateOrderDraftAction) GetInputType() reflect.Type {
	return reflect.TypeOf(CreateOrderDraftRequest{})
}

func mapCartItemsToDraftLines(cart *cartEntity.Cart) []createOrder.CartLineItem {
	out := make([]createOrder.CartLineItem, 0, len(cart.Items))
	for _, it := range cart.Items {
		out = append(out, createOrder.CartLineItem{
			ProductID: it.ProductID.String(),
			Quantity:  it.Quantity,
			Price:     it.TotalPrice,
			Comment:   it.Comment,
		})
	}
	return out
}

func draftCustomerDisplayName(fullName *string, phone string) string {
	if fullName != nil {
		n := strings.TrimSpace(*fullName)
		if n != "" {
			return n
		}
	}
	return phone
}
