package iiko

import (
	"net/http"

	apiutils "github.com/scrumno/scrumno-api/internal/api/utils"
	iikoClient "github.com/scrumno/scrumno-api/internal/iiko"
)

type GetNomenclatureAction struct {
	client iikoClient.Client
}

func NewGetNomenclatureAction(client iikoClient.Client) *GetNomenclatureAction {
	return &GetNomenclatureAction{
		client: client,
	}
}

type nomenclatureErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}

func (a *GetNomenclatureAction) Action(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("organizationId")
	if orgID == "" {
		apiutils.JSONResponse(w, nomenclatureErrorResponse{
			IsSuccess: false,
			Error:     "organizationId query parameter is required",
		}, http.StatusBadRequest)
		return
	}

	result, err := a.client.GetNomenclature(r.Context(), orgID)
	if err != nil {
		status := http.StatusBadGateway
		if result != nil && result.StatusCode >= 500 {
			status = http.StatusBadGateway
		}

		apiutils.JSONResponse(w, nomenclatureErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	_, _ = w.Write(result.Body)
}

