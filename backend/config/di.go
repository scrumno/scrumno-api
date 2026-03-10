package config

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action"
	iikoAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/iiko"
	healthAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
	checkStatusConnectDB "github.com/scrumno/scrumno-api/internal/health/query/check-status-connect-db"
	createUser "github.com/scrumno/scrumno-api/internal/users/command/create-user"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
	internalIiko "github.com/scrumno/scrumno-api/internal/iiko"
	"github.com/scrumno/scrumno-api/shared/factory"
)

func DI(cfg *Config) *action.Actions {
	// repository
	statusRepo := status.NewStatusRepository(DB)
	userRepo := factory.NewGormRepository[user.User](DB)

	// service
	checkStatusFetcher := checkStatusConnectDB.NewFetcher(statusRepo)

	// command
	createUserHandler := createUser.NewCreateUserHandler(userRepo)

	// query

	// external clients
	iikoClient := internalIiko.NewClient(cfg.Iiko)

	return &action.Actions{
		CheckStatusConnectDB: healthAction.NewCheckStatusConnectDBAction(checkStatusFetcher),

		// users
		CreateUser: userAction.NewCreateUserAction(createUserHandler),

		// iiko
		CreateIikoPickupOrder: iikoAction.NewCreatePickupOrderAction(iikoClient),
		GetIikoOrganizations:  iikoAction.NewGetOrganizationsAction(iikoClient),
		GetIikoNomenclature:   iikoAction.NewGetNomenclatureAction(iikoClient),
	}
}
