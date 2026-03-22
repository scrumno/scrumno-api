package iiko

import (
	"context"
	"fmt"
	"net/http"
	"time"

	setAccess "github.com/scrumno/scrumno-api/internal/iiko/command/auth/set-access"
	addCustomer "github.com/scrumno/scrumno-api/internal/iiko/command/order/add-customer"
	addOrderItems "github.com/scrumno/scrumno-api/internal/iiko/command/order/add-order-items"
	addOrderPayments "github.com/scrumno/scrumno-api/internal/iiko/command/order/add-order-payments"
	cancelOrder "github.com/scrumno/scrumno-api/internal/iiko/command/order/cancel-order"
	changeOrderPayments "github.com/scrumno/scrumno-api/internal/iiko/command/order/change-order-payments"
	closeOrder "github.com/scrumno/scrumno-api/internal/iiko/command/order/close-order"
	createOrder "github.com/scrumno/scrumno-api/internal/iiko/command/order/create-order"
	initByPosOrder "github.com/scrumno/scrumno-api/internal/iiko/command/order/init-by-pos-order"
	iikoconfig "github.com/scrumno/scrumno-api/internal/iiko/config"
	"github.com/scrumno/scrumno-api/internal/iiko/entity/access"
	commandStatus "github.com/scrumno/scrumno-api/internal/iiko/entity/command-status"
	externalMenu "github.com/scrumno/scrumno-api/internal/iiko/entity/external-menu"
	"github.com/scrumno/scrumno-api/internal/iiko/entity/menu"
	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
	iikoMiddleware "github.com/scrumno/scrumno-api/internal/iiko/middleware"
	getCommandStatus "github.com/scrumno/scrumno-api/internal/iiko/query/command/get-status"
	getExternalMenuByID "github.com/scrumno/scrumno-api/internal/iiko/query/external-menu/get-by-id"
	getExternalMenuList "github.com/scrumno/scrumno-api/internal/iiko/query/external-menu/get-list"
	getMenu "github.com/scrumno/scrumno-api/internal/iiko/query/menu/get-menu"
	getByIDs "github.com/scrumno/scrumno-api/internal/iiko/query/order/get-by-ids"
	getByPosOrderIDs "github.com/scrumno/scrumno-api/internal/iiko/query/order/get-by-pos-order-ids"
	authorizeService "github.com/scrumno/scrumno-api/internal/iiko/service/authorize-service"
	tokenProvider "github.com/scrumno/scrumno-api/internal/iiko/service/token-provider"
)

type Container struct {
	SetAccess           *setAccess.Handler
	CreateOrder         *createOrder.Handler
	AddCustomer         *addCustomer.Handler
	AddOrderItems       *addOrderItems.Handler
	AddOrderPayments    *addOrderPayments.Handler
	ChangeOrderPayments *changeOrderPayments.Handler
	CloseOrder          *closeOrder.Handler
	CancelOrder         *cancelOrder.Handler
	InitByPosOrder      *initByPosOrder.Handler
	GetCommandStatus    *getCommandStatus.Fetcher
	GetMenu             *getMenu.Fetcher
	GetByIDs            *getByIDs.Fetcher
	GetByPosOrderIDs    *getByPosOrderIDs.Fetcher
	GetExternalMenuList *getExternalMenuList.Fetcher
	GetExternalMenuByID *getExternalMenuByID.Fetcher
}

func NewContainer(cfg *iikoconfig.Config) *Container {
	if cfg == nil {
		c := iikoconfig.Load()
		cfg = &c
	}

	var accessRepo *access.AccessRepository
	provider := tokenProvider.NewProvider(func(ctx context.Context) (string, error) {
		if accessRepo == nil {
			return "", fmt.Errorf("репозиторий access iiko не инициализирован")
		}
		token, err := accessRepo.PostAccessToken(ctx, cfg.Login, cfg.Password)
		if err != nil {
			return "", err
		}
		return token.Token, nil
	})

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: iikoMiddleware.NewAuthRetryTransport(
			http.DefaultTransport,
			provider,
		),
	}

	// repository
	accessRepo = access.NewAccessRepository(cfg.BaseURL, httpClient)
	menuRepo := menu.NewRepository(cfg.BaseURL, httpClient)
	externalMenuRepo := externalMenu.NewRepository(cfg.BaseURL, httpClient)
	orderRepo := order.NewRepository(cfg.BaseURL, httpClient)
	commandStatusRepo := commandStatus.NewRepository(cfg.BaseURL, httpClient)

	// service
	authSvc := authorizeService.NewService(accessRepo)

	// command
	setAccessHandler := setAccess.NewHandler(authSvc, cfg.Login, cfg.Password)
	createOrderHandler := createOrder.NewHandler(orderRepo)
	addCustomerHandler := addCustomer.NewHandler(orderRepo)
	addOrderItemsHandler := addOrderItems.NewHandler(orderRepo)
	addOrderPaymentsHandler := addOrderPayments.NewHandler(orderRepo)
	changeOrderPaymentsHandler := changeOrderPayments.NewHandler(orderRepo)
	closeOrderHandler := closeOrder.NewHandler(orderRepo)
	cancelOrderHandler := cancelOrder.NewHandler(orderRepo)
	initByPosOrderHandler := initByPosOrder.NewHandler(orderRepo)
	getCommandStatusFetcher := getCommandStatus.NewFetcher(commandStatusRepo)
	menuFetcher := getMenu.NewFetcher(menuRepo)
	getByIDsFetcher := getByIDs.NewFetcher(orderRepo)
	getByPosOrderIDsFetcher := getByPosOrderIDs.NewFetcher(orderRepo)
	externalMenuListFetcher := getExternalMenuList.NewFetcher(externalMenuRepo)
	externalMenuByIDFetcher := getExternalMenuByID.NewFetcher(externalMenuRepo)

	return &Container{
		SetAccess:           setAccessHandler,
		CreateOrder:         createOrderHandler,
		AddCustomer:         addCustomerHandler,
		AddOrderItems:       addOrderItemsHandler,
		AddOrderPayments:    addOrderPaymentsHandler,
		ChangeOrderPayments: changeOrderPaymentsHandler,
		CloseOrder:          closeOrderHandler,
		CancelOrder:         cancelOrderHandler,
		InitByPosOrder:      initByPosOrderHandler,
		GetCommandStatus:    getCommandStatusFetcher,
		GetMenu:             menuFetcher,
		GetByIDs:            getByIDsFetcher,
		GetByPosOrderIDs:    getByPosOrderIDsFetcher,
		GetExternalMenuList: externalMenuListFetcher,
		GetExternalMenuByID: externalMenuByIDFetcher,
	}
}
