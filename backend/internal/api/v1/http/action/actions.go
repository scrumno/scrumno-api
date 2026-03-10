package action

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	iikoAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/iiko"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
)

type Actions struct {
	// db
	CheckStatusConnectDB *health.CheckStatusConnectDBAction

	// users
	CreateUser *userAction.CreateUserAction

	// iiko
	CreateIikoPickupOrder *iikoAction.CreatePickupOrderAction
	GetIikoOrganizations  *iikoAction.GetOrganizationsAction
	GetIikoNomenclature   *iikoAction.GetNomenclatureAction
}
