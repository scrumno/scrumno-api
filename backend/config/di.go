package config

import (
	"time"

	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	iikoMenuService "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/menu/service"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	authAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	cartAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/cart"
	cartProductAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/product"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/menu"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/orders"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
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
	refreshMenu "github.com/scrumno/scrumno-api/internal/menu/command/refresh-menu"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
	saveProductCommand "github.com/scrumno/scrumno-api/internal/products/command/save-product"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
	saveProductListener "github.com/scrumno/scrumno-api/internal/products/listener/save-product"
	updateUserProfile "github.com/scrumno/scrumno-api/internal/users/command/update-user-profile"
	conditionsUpdateProfilePolicy "github.com/scrumno/scrumno-api/internal/users/service/conditions-update-profile"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
	"github.com/scrumno/scrumno-api/shared/services/jwt"
	"github.com/scrumno/scrumno-api/shared/services/sms"

	// Cart
	addProduct "github.com/scrumno/scrumno-api/internal/cart/command/add-product-to-cart"
	clearCart "github.com/scrumno/scrumno-api/internal/cart/command/clear-cart"
	createCart "github.com/scrumno/scrumno-api/internal/cart/command/create-cart"
	removeProduct "github.com/scrumno/scrumno-api/internal/cart/command/remove-product"
	updateProduct "github.com/scrumno/scrumno-api/internal/cart/command/update-product-cart"
	cart "github.com/scrumno/scrumno-api/internal/cart/entity"
	getCart "github.com/scrumno/scrumno-api/internal/cart/query/get-cart-by-user-id"
	"github.com/scrumno/scrumno-api/shared/services/snapshot"
	fileStorage "github.com/scrumno/scrumno-api/shared/services/storage"
)

func DI() (*action.Actions, *action.Listeners) {
	cfg := Load()

	em := eventManager.New()

	/* INTEGRATION SYSTEMs */

	var (
		// service
		orderProvider   interfaces.OrderProvider
		orderBuilder    interfaces.OrderBuilder
		snapshotService interfaces.SnapshotService
		snapshotStore   interfaces.SnapshotStore

		menuProvider interfaces.MenuProvider

		// handlers
		getMenuHandler interfaces.GetMenuHandler

		// actions
		refreshMenuAction menu.RefreshMenuAction

		// config
	)
	var iikoCfg *iikoConfig.Config

	switch cfg.IntegrationSystem.IntegrationSystem {

	case "iiko":
		iikoCfg = iikoConfig.Load()

		// services
		menuProvider = iikoMenuService.NewMenuProvider(iikoCfg)
		snapshotStore = fileStorage.NewFileStore(iikoCfg.SnapshotFilePath)
		snapshotService = snapshot.NewSnapshotService(snapshotStore)

		// handlers
		getMenuHandler = refreshMenu.NewHandler(menuProvider, em, snapshotService)

		// actions
		refreshMenuAction = menu.NewRefreshMenuAction(getMenuHandler)
	}

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
	productRepo := product.NewProductRepository(DB)
	cartRepo := cart.NewCartRepository(DB)

	jwtManager := jwt.NewManager(jwt.Config{
		AccessSecret:    string(cfg.JWT.SecretKey),
		RefreshSecret:   string(cfg.JWT.SecretKey),
		AccessTokenTtl:  15 * time.Minute,
		RefreshTokenTtl: 7 * 24 * time.Hour,
	})

	// service
	checkStatusFetcher := checkStatusConnectDb.NewFetcher(statusRepo)

	// service (нужен до createAuthorizeCodeHandler)
	createUniqueCodeSvc := createUniqueCode.NewCreateUniqueCodeService()

	// command
	conditionsProfilePolicy := conditionsUpdateProfilePolicy.NewHandler()
	updateUserProfileHandler := updateUserProfile.NewHandler(registrationRepo, conditionsProfilePolicy)
	logoutHandler := logout.NewHandler(tokensRepo)
	checkOntimeCodeHandler := checkOntimeCode.NewHandler(codesRepo)
	createUserAuthHandler := createUserAuth.NewHandler(registrationRepo)
	createAuthorizeTokensHandler := createAuthorizeTokens.NewHandler(tokensRepo, jwtManager)
	createAuthorizeCodeHandler := createAuthorizeCode.NewHandler(codesRepo, createUniqueCodeSvc)

	createOrderHandler := createOrder.NewHandler(orderProvider, orderBuilder)

	saveProductHandler := saveProductCommand.NewHandler(productRepo)

	createCartHandler := createCart.NewHandler(cartRepo)
	clearCartHandler := clearCart.NewHandler(cartRepo)
	addProductHandler := addProduct.NewHandler(cartRepo)
	removeProductHandler := removeProduct.NewHandler(cartRepo)
	updateProductHandler := updateProduct.NewHandler(cartRepo)
	getCartFetcher := getCart.NewFetcher(cartRepo)

	// query
	getRefreshTokensFetcher := getRefreshTokensAvailable.NewFetcher(tokensRepo, jwtManager)
	findUserByPhoneFetcher := findUserByPhone.NewFetcher(registrationRepo)
	getSmsCodeSendAvailableFetcher := getSmsCodeSendAvailable.NewFetcher(codesRepo)
	getSmsCodeFetcher := getSmsCode.NewFetcher(smsService)

	// listeners
	saveProductListener := saveProductListener.NewListener(saveProductHandler)

	return &action.Actions{
			CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher),

			// users
			UpdateUserProfile: user.NewUpdateUserProfileAction(updateUserProfileHandler),

			// auth
			Registration:  authAction.NewRegistrationAction(findUserByPhoneFetcher, checkOntimeCodeHandler, createUserAuthHandler, createAuthorizeTokensHandler),
			Authorize:     authAction.NewAuthorizeAction(findUserByPhoneFetcher, checkOntimeCodeHandler, createAuthorizeTokensHandler),
			Logout:        authAction.NewLogoutAction(logoutHandler, findUserByPhoneFetcher),
			RefreshTokens: authAction.NewRefreshTokensAction(getRefreshTokensFetcher, findUserByPhoneFetcher, createAuthorizeTokensHandler),

			JWTManager: jwtManager,
			SmsCode:    authAction.NewAuthCodeAction(getSmsCodeSendAvailableFetcher, getSmsCodeFetcher, createAuthorizeCodeHandler),

			// orders
			CreateOrder: orders.NewCreateOrderAction(createOrderHandler),

			// cart
			CreateCart: cartAction.NewCreateAction(createCartHandler),
			ClearCart:  cartAction.NewClearAction(clearCartHandler),
			GetCart:    cartAction.NewGetCartAction(getCartFetcher),

			AddProductToCart:      cartProductAction.NewAddProductAction(addProductHandler),
			RemoveProductFromCart: cartProductAction.NewRemoveProductAction(removeProductHandler),
			UpdateProductFromCart: cartProductAction.NewUpdateAction(updateProductHandler),

			// общие экшены для всех интеграционных систем
			RefreshMenu: &refreshMenuAction,
		},
		&action.Listeners{
			SaveProduct: saveProductListener,
		}
}
