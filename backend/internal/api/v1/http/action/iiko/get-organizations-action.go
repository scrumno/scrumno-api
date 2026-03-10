package iiko

import (
	"net/http"

	apiutils "github.com/scrumno/scrumno-api/internal/api/utils"
	iikoClient "github.com/scrumno/scrumno-api/internal/iiko"
)

type GetOrganizationsAction struct {
	client iikoClient.Client
}

func NewGetOrganizationsAction(client iikoClient.Client) *GetOrganizationsAction {
	return &GetOrganizationsAction{
		client: client,
	}
}

type organizationsErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}

func (a *GetOrganizationsAction) Action(w http.ResponseWriter, r *http.Request) {
	result, err := a.client.GetOrganizations(r.Context())
	if err != nil {
		status := http.StatusBadGateway
		if result != nil && result.StatusCode >= 500 {
			status = http.StatusBadGateway
		}

		apiutils.JSONResponse(w, organizationsErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	_, _ = w.Write(result.Body)
}

