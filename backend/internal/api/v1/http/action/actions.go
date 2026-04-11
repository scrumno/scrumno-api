package action

import (
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/auth"
	cartAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/cart"
	cartProductAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/product"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/health"
	menuAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/menu"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/orders"
	queueAction "github.com/scrumno/scrumno-api/internal/api/v1/http/action/queue"
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

	// cart
	CreateCart            *cartAction.CreateAction
	ClearCart             *cartAction.ClearAction
	AddProductToCart      *cartProductAction.AddProductAction
	RemoveProductFromCart *cartProductAction.RemoveProductAction
	UpdateProductFromCart *cartProductAction.UpdateAction
	GetCart               *cartAction.GetCartAction

	// iiko
	RefreshMenu *menuAction.RefreshMenuAction
	RefreshQueue *queueAction.RefreshQueueAction

	// queue
	GetQueue         *queueAction.GetQueueAction
	AddInQueue       *queueAction.AddInQueueAction
	GetNearestRange  *queueAction.GetNearestRangeAction

	// auth
	Registration  *auth.RegistrationAction
	Authorize     *auth.AuthorizeAction
	RefreshTokens *auth.RefreshTokensAction
	Logout        *auth.LogoutAction

	GetMenu *menuAction.GetMenuAction

	JWTManager *jwt.Manager

	SmsCode *auth.AuthCodeAction
}

type Listeners struct {
	SaveProduct  *saveProduct.Listener
	SaveModifier *saveModifier.Listener
	SaveMenu     *saveMenu.Listener
}
