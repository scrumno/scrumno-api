package config

import (
	"time"

	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	authAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	logout "github.com/scrumno/scrumno-api/internal/authorize/command/logout"
	getRefreshTokensAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-refresh-tokens-available"
	createUniqueCode "github.com/scrumno/scrumno-api/internal/authorize/service/create-unique-code"
	authEntity "github.com/scrumno/scrumno-api/internal/authorize/entity"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
    tokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDB "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
	userEntity "github.com/scrumno/scrumno-api/internal/users/entity/user"
	"github.com/scrumno/scrumno-api/shared/factory"
	"github.com/scrumno/scrumno-api/shared/jwt"
	findUSerByPhoneQuery "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	"github.com/scrumno/scrumno-api/shared/sms"
	checkOnetimeCodeCommand "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createUserCommand "github.com/scrumno/scrumno-api/internal/authorize/command/create-user"
	createAuthorizeTokensCommand "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
	getSmsCodeSendAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code-send-available"
	getSmsCode "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code"
	createAuthorizeCodeCommand "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-code"
)

func DI() *action.Actions {
	cfg := Load()

	smsService := sms.NewSmsService(sms.Config{
		ApiKey: cfg.Sms.ApiKey,
		ApiPhoneNumber: cfg.Sms.ApiPhoneNumber,
	})

	// repository
	statusRepo := status.NewStatusRepository(DB)
	userRepo := factory.NewGormRepository[userEntity.User](DB)
	registrationRepo := authEntity.NewRegistrationRepository(DB)
	tokensRepo := tokens.NewTokensRepository(DB)
	codesRepo := codes.NewSmsCodesRepository(DB)

	jwtManager := jwt.NewManager(jwt.Config{
		AccessSecret:    string(cfg.JWT.SecretKey),
		RefreshSecret:   string(cfg.JWT.SecretKey),
		AccessTokenTtl:  15 * time.Minute,
		RefreshTokenTtl: 7 * 24 * time.Hour,
	})

	// service
	checkStatusFetcher := checkStatusConnectDB.NewFetcher(statusRepo)

	// service (нужен до createAuthorizeCodeHandler)
	createUniqueCodeSvc := createUniqueCode.NewCreateUniqueCodeService()

	// command
	createUserHandler := createUser.NewCreateUserHandler(userRepo)
	logoutHandler := logout.NewHandler(tokensRepo)
	checkOnetimeCodeHandler := checkOnetimeCodeCommand.NewHandler(codesRepo)
	createUserCommandHandler := createUserCommand.NewHandler(registrationRepo)
	createAuthorizeTokensHandler := createAuthorizeTokensCommand.NewHandler(tokensRepo, jwtManager)
	createAuthorizeCodeHandler := createAuthorizeCodeCommand.NewHandler(codesRepo, createUniqueCodeSvc)

	// query
	getRefreshTokensFetcher := getRefreshTokensAvailable.NewFetcher(tokensRepo, jwtManager)
	findUserByPhoneFetcher := findUSerByPhoneQuery.NewFetcher(registrationRepo)
	getSmsCodeSendAvailableFetcher := getSmsCodeSendAvailable.NewFetcher(codesRepo)
	getSmsCodeFetcher := getSmsCode.NewFetcher(smsService)

	return &action.Actions{
		CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher),

		// users
		CreateUser: userAction.NewCreateUserAction(createUserHandler),

		// auth
		Registration: authAction.NewRegistrationAction(findUserByPhoneFetcher, checkOnetimeCodeHandler, createUserCommandHandler, createAuthorizeTokensHandler),
		Authorize:     authAction.NewAuthorizeAction(findUserByPhoneFetcher, checkOnetimeCodeHandler, createAuthorizeTokensHandler),
		Logout:        authAction.NewLogoutAction(logoutHandler, findUserByPhoneFetcher),
		RefreshTokens: authAction.NewRefreshTokensAction(getRefreshTokensFetcher, findUserByPhoneFetcher, createAuthorizeTokensHandler),

		JWTManager: jwtManager,
		SmsCode:    authAction.NewAuthCodeAction(getSmsCodeSendAvailableFetcher, getSmsCodeFetcher, createAuthorizeCodeHandler),
	}
}
