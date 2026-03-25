package config

import (
	"time"

	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	authAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/orders"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	checkOntimeCode "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createAuthorizeCode "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-code"
	createAuthorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
	createUserAuth "github.com/scrumno/scrumno-api/internal/authorize/command/create-user"
	logout "github.com/scrumno/scrumno-api/internal/authorize/command/logout"
	authEntity "github.com/scrumno/scrumno-api/internal/authorize/entity"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	authorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/entity/tokens"
	findUserByPhone "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	getRefreshTokensAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-refresh-tokens-available"
	getSmsCode "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code"
	getSmsCodeSendAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code-send-available"
	createUniqueCode "github.com/scrumno/scrumno-api/internal/authorize/service/create-unique-code"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDb "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
	userEntity "github.com/scrumno/scrumno-api/internal/users/entity/user"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/services/jwt"
	"github.com/scrumno/scrumno-api/shared/services/sms"
)

func DI(cfg *Config) *action.Actions {

	/* INTEGRATION SYSTEMs */

	var (
		// service
		orderProvider interfaces.OrderProvider
		orderBuilder  interfaces.OrderBuilder

		// config

	)

	/* INTEGRATION SYSTEMs END */

	smsService := sms.NewSmsService(sms.Config{
		ApiKey:         cfg.Sms.ApiKey,
		ApiPhoneNumber: cfg.Sms.ApiPhoneNumber,
	})
	// repository
	statusRepo := status.NewStatusRepository(DB)
	registrationRepo := authEntity.NewRegistrationRepository(DB)
	tokensRepo := authorizeTokens.NewTokensRepository(DB)
	codesRepo := codes.NewSmsCodesRepository(DB)

	jwtManager := jwt.NewManager(jwt.Config{
		AccessSecret: cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessTokenTtl: cfg.JWT.AccessTokenTtl,
		RefreshTokenTtl: cfg.JWT.RefreshTokenTtl,
	})

	// service
	checkStatusFetcher := checkStatusConnectDb.NewFetcher(statusRepo)

	// service (нужен до createAuthorizeCodeHandler)
	createUniqueCodeSvc := createUniqueCode.NewCreateUniqueCodeService()

	// command
	conditionsUpdateProfilePolicy := conditionsUpdateProfilePolicy.NewHandler()
	updateUserProfileHandler := updateUserProfile.NewHandler(registrationRepo, conditionsUpdateProfilePolicy)
	logoutHandler := logout.NewHandler(tokensRepo)
	checkOntimeCodeHandler := checkOntimeCode.NewHandler(codesRepo)
	createUserAuthHandler := createUserAuth.NewHandler(registrationRepo)
	createAuthorizeTokensHandler := createAuthorizeTokens.NewHandler(tokensRepo, jwtManager)
	createAuthorizeCodeHandler := createAuthorizeCode.NewHandler(codesRepo, createUniqueCodeSvc)

	createOrderHandler := createOrder.NewHandler(orderProvider, orderBuilder)
	// query
	getRefreshTokensFetcher := getRefreshTokensAvailable.NewFetcher(tokensRepo, jwtManager)
	findUserByPhoneFetcher := findUserByPhone.NewFetcher(registrationRepo)
	getSmsCodeSendAvailableFetcher := getSmsCodeSendAvailable.NewFetcher(codesRepo)
	getSmsCodeFetcher := getSmsCode.NewFetcher(smsService)

	return &action.Actions{
		CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher),

		// users
		UpdateUserProfile: userAction.NewUpdateUserProfileAction(updateUserProfileHandler),

		// auth
		Registration:  authAction.NewRegistrationAction(findUserByPhoneFetcher, checkOntimeCodeHandler, createUserAuthHandler, createAuthorizeTokensHandler),
		Authorize:     authAction.NewAuthorizeAction(findUserByPhoneFetcher, checkOntimeCodeHandler, createAuthorizeTokensHandler),
		Logout:        authAction.NewLogoutAction(logoutHandler, findUserByPhoneFetcher),
		RefreshTokens: authAction.NewRefreshTokensAction(getRefreshTokensFetcher, findUserByPhoneFetcher, createAuthorizeTokensHandler),

		JWTManager: jwtManager,
		SmsCode:    authAction.NewAuthCodeAction(getSmsCodeSendAvailableFetcher, getSmsCodeFetcher, createAuthorizeCodeHandler),

		// orders
		CreateOrder: orders.NewCreateOrderAction(createOrderHandler),
	}
}
