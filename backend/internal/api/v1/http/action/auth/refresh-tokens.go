package auth

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	getRefreshTokensFetcher "github.com/scrumno/scrumno-api/internal/authorize/query/get-refresh-tokens-available"
	findUserByPhoneFetcher "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	createAuthorizeTokensHandler "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
)

type RefreshTokensAction struct {
	GetRefreshTokensFetcher *getRefreshTokensFetcher.Fetcher
	FindUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher
	CreateAuthorizeTokensHandler *createAuthorizeTokensHandler.Handler
}

func NewRefreshTokensAction(
	getRefreshTokensFetcher *getRefreshTokensFetcher.Fetcher, 
	findUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher, 
	createAuthorizeTokensHandler *createAuthorizeTokensHandler.Handler,
) *RefreshTokensAction {
	return &RefreshTokensAction{
		GetRefreshTokensFetcher: getRefreshTokensFetcher,
		FindUserByPhoneFetcher: findUserByPhoneFetcher,
		CreateAuthorizeTokensHandler: createAuthorizeTokensHandler,
	}
}

func (a *RefreshTokensAction) GetInputType() reflect.Type {
	return reflect.TypeOf(RefreshTokensRequest{})
}

type RefreshTokensRequest struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2VmN2Q0OTgtYmU4Yy00ZGRkLTljM2UtMTY2MGFjODVkYWIxIiwicGhvbmUiOiI3OTE0NzgxNzgzMiIsInV1aWQiOiIiLCJzZXNzaW9uX2lkIjoiNTc5YWUyZjQtZDhjMS00YTdkLWFlN2QtZTM2YmQzNmI2YjBiIiwic3ViIjoiN2VmN2Q0OTgtYmU4Yy00ZGRkLTljM2UtMTY2MGFjODVkYWIxIiwiZXhwIjoxNzc0MDcwNDUyLCJpYXQiOjE3NzM0NjU2NTJ9.Qh3PL4RooYM9QxJfhREiRC0C2xDsSZb6v_hZWC5FN4U"`
}

func (a *RefreshTokensAction) Action(w http.ResponseWriter, r *http.Request) {

	var req RefreshTokensRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, RefreshTokensErrorResponse{
			IsSuccess: false, 
			Error: "Неверный формат запроса",
		}, http.StatusBadRequest)
		return
	}
	
	if req.RefreshToken == "" {
		utils.JSONResponse(w, RefreshTokensErrorResponse{
			IsSuccess: false, 
			Error: "Рефреш токен не может быть пустым",
		}, http.StatusBadRequest)
		return
	}

	claims, err := a.GetRefreshTokensFetcher.Fetch(r.Context(), req.RefreshToken)
	if err != nil {
		utils.JSONResponse(w, RefreshTokensErrorResponse{
			IsSuccess: false,
			Error: err.Error(),
		}, http.StatusUnauthorized)
		return
	}

	user, err := a.FindUserByPhoneFetcher.Fetch(r.Context(), claims.Phone)
	if err != nil {
		utils.JSONResponse(w, RefreshTokensErrorResponse{
			IsSuccess: false,
			Error: err.Error(),
		}, http.StatusUnauthorized)
		return
	}
	if user == nil {
		utils.JSONResponse(w, RefreshTokensErrorResponse{
			IsSuccess: false,
			Error:     "Пользователь не найден",
		}, http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, expiresIn, err := a.CreateAuthorizeTokensHandler.Handle(r.Context(), createAuthorizeTokensHandler.Command{
		Phone: user.Phone,
		UserID: user.ID,
		SessionID: claims.SessionID,
		RevokePreviousToken: true,
	})
	if err != nil {
		utils.JSONResponse(w, RefreshTokensErrorResponse{
			IsSuccess: false,
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, RefreshTokensResponse{
		IsSuccess:    true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, http.StatusOK)
}

type RefreshTokensResponse struct {
	IsSuccess    bool   `json:"isSuccess"`
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    ExpiresIn    int64  `json:"expiresIn"`
}

type RefreshTokensErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}