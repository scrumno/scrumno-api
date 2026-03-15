package auth

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	findUserByPhoneFetcher "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	checkOnetimeCodeHandler "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createAuthorizeTokensHandler "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
)

type AuthorizeAction struct {
	FindUserByPhoneFetcher       *findUserByPhoneFetcher.Fetcher
	CheckOnetimeCodeHandler     *checkOnetimeCodeHandler.Handler
	CreateAuthorizeTokensHandler *createAuthorizeTokensHandler.Handler
}

func NewAuthorizeAction(
	findUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher,
	checkOnetimeCodeHandler *checkOnetimeCodeHandler.Handler,
	createAuthorizeTokensHandler *createAuthorizeTokensHandler.Handler,
) *AuthorizeAction {
	return &AuthorizeAction{
		FindUserByPhoneFetcher: findUserByPhoneFetcher,
		CheckOnetimeCodeHandler: checkOnetimeCodeHandler,
		CreateAuthorizeTokensHandler: createAuthorizeTokensHandler,
	}
}	

func (a *AuthorizeAction) GetInputType() reflect.Type {
    return reflect.TypeOf(AuthorizeRequest{})
}

type AuthorizeRequest struct {
	Phone string `json:"phone" example:"79099000000"`
	Code  string `json:"code" example:"1234"`
}

func (a *AuthorizeAction) Action(w http.ResponseWriter, r *http.Request) {

	var req AuthorizeRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, AuthorizeErrorResponse{
			IsSuccess: false, 
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	user, err := a.FindUserByPhoneFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, AuthorizeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if user == nil {
		utils.JSONResponse(w, AuthorizeErrorResponse{
			IsSuccess: false,
			Error:     "Пользователь не найден",
		}, http.StatusBadRequest)
		return
	}

	checkCmd := checkOnetimeCodeHandler.Command{
		Phone:    req.Phone,
		Code:     req.Code,
		CodeType: codes.AuthType,
	}

	err = a.CheckOnetimeCodeHandler.Handle(r.Context(), checkCmd)
	if err != nil {
		utils.JSONResponse(w, AuthorizeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	tokensCmd := createAuthorizeTokensHandler.Command{
		Phone:  user.Phone,
		UserID: user.ID,
		SessionID: "",
		RevokePreviousToken: false,
	}

	accessToken, refreshToken, expiresIn, err := a.CreateAuthorizeTokensHandler.Handle(r.Context(), tokensCmd)
	if err != nil {
		utils.JSONResponse(w, AuthorizeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, AuthorizeResponse{
		IsSuccess:    true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, http.StatusOK)
}

type AuthorizeResponse struct {
    IsSuccess    bool   `json:"isSuccess"`
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    ExpiresIn    int64  `json:"expiresIn"`
}

type AuthorizeErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}

