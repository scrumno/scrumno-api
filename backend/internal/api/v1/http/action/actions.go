package action

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/orders"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	"github.com/scrumno/scrumno-api/shared/services/jwt"
)

type Actions struct {
	// db
	CheckStatusConnectDB *health.CheckStatusConnectDBAction

	// users
	CreateUser *userAction.CreateUserAction

	// orders
	CreateOrder *orders.CreateOrderAction

	// auth
	Registration  *auth.RegistrationAction
	Authorize     *auth.AuthorizeAction
	RefreshTokens *auth.RefreshTokensAction
	Logout        *auth.LogoutAction

	JWTManager *jwt.Manager

	SmsCode *auth.AuthCodeAction
}
