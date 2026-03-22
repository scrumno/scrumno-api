package action

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	userAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/user"
	setAccess "github.com/scrumno/scrumno-api/internal/iiko/command/auth/set-access"
	"github.com/scrumno/scrumno-api/shared/jwt"
)

type Actions struct {
	// db
	CheckStatusConnectDB *health.CheckStatusConnectDBAction

	// users
	UpdateUserProfile *userAction.UpdateUserProfileAction

	// auth
	Registration  *auth.RegistrationAction
	Authorize     *auth.AuthorizeAction
	RefreshTokens *auth.RefreshTokensAction
	Logout        *auth.LogoutAction

	JWTManager *jwt.Manager

	SmsCode *auth.AuthCodeAction

	// iiko
	SetAccess *setAccess.Handler
}
