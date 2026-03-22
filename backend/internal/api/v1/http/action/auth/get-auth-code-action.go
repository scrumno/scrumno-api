package auth

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	createAuthorizeCode "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-code"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	getSmsCode "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code"
	getSmsCodeSendAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code-send-available"
)

type AuthCodeAction struct {
	GetSmsCodeSendAvailableFetcher *getSmsCodeSendAvailable.Fetcher
	GetSmsCodeFetcher              *getSmsCode.Fetcher
	CreateAuthorizeCodeHandler     *createAuthorizeCode.Handler
}

func NewAuthCodeAction(
	getSmsCodeSendAvailableFetcher *getSmsCodeSendAvailable.Fetcher,
	getSmsCodeFetcher *getSmsCode.Fetcher,
	createAuthorizeCodeHandler *createAuthorizeCode.Handler,
) *AuthCodeAction {
	return &AuthCodeAction{
		GetSmsCodeSendAvailableFetcher: getSmsCodeSendAvailableFetcher,
		GetSmsCodeFetcher:              getSmsCodeFetcher,
		CreateAuthorizeCodeHandler:     createAuthorizeCodeHandler,
	}
}

func (a *AuthCodeAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetAuthCodeRequest{})
}

type GetAuthCodeRequest struct {
	Phone    string          `json:"phone" example:"79090000000"`
	CodeType codes.CodesType `json:"codeType" example:"authorize"`
}

func (a *AuthCodeAction) Action(w http.ResponseWriter, r *http.Request) {

	var req GetAuthCodeRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Неверный формат запроса",
		}, http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Укажите номер телефона",
		}, http.StatusBadRequest)
		return
	}

	if req.CodeType != codes.AuthType && req.CodeType != codes.RegisterType {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Недопустимый тип кода",
		}, http.StatusBadRequest)
		return
	}

	_, err = a.GetSmsCodeSendAvailableFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	authorizeCode, err := a.CreateAuthorizeCodeHandler.Handle(
		r.Context(),
		req.Phone,
		req.CodeType,
	)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("Ваш код авторизации: %s", authorizeCode.Code)

	_, err = a.GetSmsCodeFetcher.Fetch(r.Context(), req.Phone, message)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, AuthCodeResponse{
		IsSuccess: true,
	}, http.StatusOK)
}

type AuthCodeResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

type AuthCodeErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}
