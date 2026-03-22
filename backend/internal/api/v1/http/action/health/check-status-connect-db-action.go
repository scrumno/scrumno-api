package health

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	checkStatusConnectDb "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
)

type Request struct {
	Phone    string `json:"phone" example:"79099000000"`
	FullName string `json:"full_name" example:"Иван Аресньев"`
}

func (a *CheckStatusConnectDBAction) GetInputType() reflect.Type {
	return reflect.TypeOf(Request{})
}

type CheckStatusConnectDBAction struct {
	fetcher *checkStatusConnectDb.Fetcher
}

func NewCheckStatusConnectDBAction(fetcher *checkStatusConnectDb.Fetcher) *CheckStatusConnectDBAction {
	return &CheckStatusConnectDBAction{fetcher: fetcher}
}

func (a *CheckStatusConnectDBAction) Action(w http.ResponseWriter, _ *http.Request) {
	dto := a.fetcher.Fetch(checkStatusConnectDb.Query{})

	if !dto.IsConnected {
		utils.JSONResponse(w, map[string]bool{"isOk": false}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, map[string]bool{"isOk": true}, http.StatusOK)
}
