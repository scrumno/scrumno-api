package iiko

import (
	"net/http"

	apiutils "github.com/scrumno/scrumno-api/internal/api/utils"
	iikoClient "github.com/scrumno/scrumno-api/internal/iiko"
)

type GetTerminalsAction struct {
	client iikoClient.Client
}

func NewGetTerminalsAction(client iikoClient.Client) *GetTerminalsAction {
	return &GetTerminalsAction{
		client: client,
	}
}

type terminalsErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}

func (a *GetTerminalsAction) Action(w http.ResponseWriter, r *http.Request) {
	result, err := a.client.GetTerminals(r.Context())
	if err != nil {
		status := http.StatusBadGateway
		if result != nil && result.StatusCode >= 500 {
			status = http.StatusBadGateway
		}

		apiutils.JSONResponse(w, terminalsErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	_, _ = w.Write(result.Body)
}

