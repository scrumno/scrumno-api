package iiko

import (
	"encoding/json"
	"net/http"

	apiutils "github.com/scrumno/scrumno-api/internal/api/utils"
	iikoClient "github.com/scrumno/scrumno-api/internal/iiko"
)

type CreatePickupOrderAction struct {
	client iikoClient.Client
}

func NewCreatePickupOrderAction(client iikoClient.Client) *CreatePickupOrderAction {
	return &CreatePickupOrderAction{
		client: client,
	}
}

type errorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}

type successResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Status    int    `json:"iikoStatus"`
	Body      string `json:"iikoBody"`
}

func (a *CreatePickupOrderAction) Action(w http.ResponseWriter, r *http.Request) {
	var order iikoClient.PickupOrder

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		apiutils.JSONResponse(w, errorResponse{IsSuccess: false, Error: "invalid JSON body"}, http.StatusBadRequest)
		return
	}

	if len(order.Items) == 0 {
		apiutils.JSONResponse(w, errorResponse{IsSuccess: false, Error: "items must not be empty"}, http.StatusBadRequest)
		return
	}

	if order.Customer.Phone == "" {
		apiutils.JSONResponse(w, errorResponse{IsSuccess: false, Error: "customer.phone is required"}, http.StatusBadRequest)
		return
	}

	result, err := a.client.CreatePickupOrder(r.Context(), order)
	if err != nil {
		status := http.StatusInternalServerError
		if result != nil && result.StatusCode >= 400 {
			status = http.StatusBadGateway
		}

		apiutils.JSONResponse(w, errorResponse{IsSuccess: false, Error: err.Error()}, status)
		return
	}

	body := ""
	if result != nil && len(result.Body) > 0 {
		body = string(result.Body)
	}

	apiutils.JSONResponse(w, successResponse{
		IsSuccess: true,
		Status:    result.StatusCode,
		Body:      body,
	}, http.StatusOK)
}

