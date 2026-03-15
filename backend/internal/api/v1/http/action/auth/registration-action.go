package auth

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	findUserByPhoneFetcher "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	checkOnetimeCodeHandler "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createUserHandler "github.com/scrumno/scrumno-api/internal/authorize/command/create-user"
	createAuthorizeTokensHandler "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
)

type RegistrationAction struct {
	FindUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher
	CheckOnetimeCodeHandler *checkOnetimeCodeHandler.Handler
	CreateUserHandler *createUserHandler.Handler
	CreateAuthorizeTokensHandler *createAuthorizeTokensHandler.Handler
}

func NewRegistrationAction(
	findUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher, 
	checkOnetimeCodeHandler *checkOnetimeCodeHandler.Handler,
	createUserHandler *createUserHandler.Handler,
	createAuthorizeTokensHandler *createAuthorizeTokensHandler.Handler,
) *RegistrationAction {
	return &RegistrationAction{
		FindUserByPhoneFetcher: findUserByPhoneFetcher,
		CheckOnetimeCodeHandler: checkOnetimeCodeHandler,
		CreateUserHandler: createUserHandler,
		CreateAuthorizeTokensHandler: createAuthorizeTokensHandler,
	}
}

func (a *RegistrationAction) GetInputType() reflect.Type {
    return reflect.TypeOf(RegistrationRequest{})
}

type RegistrationRequest struct {
	Phone string `json:"phone" example:"79099000000"`
	Code  string `json:"code" example:"1234"`
}

func (a *RegistrationAction) Action(w http.ResponseWriter, r *http.Request) {

	var req RegistrationRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false, 
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	user, err := a.FindUserByPhoneFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if user != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     "Не удалось создать пользователя",
		}, http.StatusBadRequest)
		return
	}

	cmd := checkOnetimeCodeHandler.Command{
		Phone: req.Phone,
		Code:  req.Code,
		CodeType: codes.RegisterType,
	}

	err = a.CheckOnetimeCodeHandler.Handle(r.Context(), cmd)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	createdUser, err := a.CreateUserHandler.Handle(
		r.Context(),
		createUserHandler.Command{
			Phone: req.Phone,
		},
	)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, expiresIn, err := a.CreateAuthorizeTokensHandler.Handle(
		r.Context(),
		createAuthorizeTokensHandler.Command{
			Phone:  createdUser.Phone,
			UserID: createdUser.ID,
		},
	)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)	
		return
	}

	utils.JSONResponse(w, RegistrationResponse{
		IsSuccess:    true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, http.StatusOK)
}

type RegistrationResponse struct {
	IsSuccess    bool   `json:"isSuccess"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type RegistrationErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}
