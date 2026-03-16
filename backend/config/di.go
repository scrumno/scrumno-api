package config

import (
	"log/slog"
	"time"

	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	authAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	iikoAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/iiko"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	checkOnetimeCodeCommand "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createAuthorizeCodeCommand "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-code"
	createAuthorizeTokensCommand "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
	createUserCommand "github.com/scrumno/scrumno-api/internal/authorize/command/create-user"
	logout "github.com/scrumno/scrumno-api/internal/authorize/command/logout"
	authEntity "github.com/scrumno/scrumno-api/internal/authorize/entity"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	tokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	findUSerByPhoneQuery "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	getRefreshTokensAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-refresh-tokens-available"
	getSmsCode "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code"
	getSmsCodeSendAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code-send-available"
	createUniqueCode "github.com/scrumno/scrumno-api/internal/authorize/service/create-unique-code"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDB "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	internalIiko "github.com/scrumno/scrumno-api/internal/iiko"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
	userEntity "github.com/scrumno/scrumno-api/internal/users/entity/user"
	"github.com/scrumno/scrumno-api/shared/factory"
	"github.com/scrumno/scrumno-api/shared/jwt"
	"github.com/scrumno/scrumno-api/shared/sms"
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

	// external clients
	slog.Info("IIKO_CONFIG",
		"baseURL", cfg.Iiko.BaseURL,
		"login", cfg.Iiko.Login,
		"password_set", cfg.Iiko.Password != "",
		"orgID", cfg.Iiko.OrganizationID,
		"terminalID", cfg.Iiko.TerminalID,
	)

	iikoClient := internalIiko.NewClient(internalIiko.Config{
		BaseURL:        cfg.Iiko.BaseURL,
		Login:          cfg.Iiko.Login,
		Password:       cfg.Iiko.Password,
		OrganizationID: cfg.Iiko.OrganizationID,
		TerminalID:     cfg.Iiko.TerminalID,
	})

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

		// iiko
		CreateIikoPickupOrder: iikoAction.NewCreatePickupOrderAction(iikoClient),
		GetIikoOrganizations:  iikoAction.NewGetOrganizationsAction(iikoClient),
		GetIikoNomenclature:   iikoAction.NewGetNomenclatureAction(iikoClient),
		GetIikoTerminals:      iikoAction.NewGetTerminalsAction(iikoClient),
	}
}
