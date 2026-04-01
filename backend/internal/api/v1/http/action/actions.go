package action

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	menuAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/menu"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/orders"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	saveMenu "github.com/scrumno/scrumno-api/internal/menu/listener/save-menu"
	saveModifier "github.com/scrumno/scrumno-api/internal/products/listener/save-modifier"
	saveProduct "github.com/scrumno/scrumno-api/internal/products/listener/save-product"
	"github.com/scrumno/scrumno-api/shared/services/jwt"
)

type Actions struct {
	// db
	CheckStatusConnectDB *health.CheckStatusConnectDBAction

	// users
	UpdateUserProfile *userAction.UpdateUserProfileAction

	// orders
	CreateOrder *orders.CreateOrderAction

	// iiko
	RefreshMenu *menuAction.RefreshMenuAction

	// auth
	Registration  *auth.RegistrationAction
	Authorize     *auth.AuthorizeAction
	RefreshTokens *auth.RefreshTokensAction
	Logout        *auth.LogoutAction

	JWTManager *jwt.Manager

	SmsCode *auth.AuthCodeAction
}

type Listeners struct {
	SaveProduct  *saveProduct.Listener
	SaveModifier *saveModifier.Listener
	SaveMenu     *saveMenu.Listener
}
