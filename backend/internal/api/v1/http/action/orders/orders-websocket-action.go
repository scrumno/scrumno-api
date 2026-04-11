package orders

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
	ordersService "github.com/scrumno/scrumno-api/internal/orders/service"
)

var ordersWebSocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type OrdersWebSocketAction struct {
	Hub  ordersService.OrdersWebSocketHub
	Repo ordersEntity.OrderRepository
}

type wsCommand struct {
	Action  string    `json:"action"`
	OrderID uuid.UUID `json:"order_id"`
}

func NewOrdersWebSocketAction(hub ordersService.OrdersWebSocketHub, repo ordersEntity.OrderRepository) *OrdersWebSocketAction {
	return &OrdersWebSocketAction{
		Hub:  hub,
		Repo: repo,
	}
}

func (a *OrdersWebSocketAction) Action(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromRequest(r)
	if claims == nil || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	conn, err := ordersWebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	connectionID := uuid.NewString()
	a.Hub.Register(connectionID, conn)
	defer func() {
		_ = a.Repo.DeactivateSubscriber(r.Context(), connectionID, nil)
		a.Hub.Unregister(connectionID)
	}()

	for {
		var cmd wsCommand
		if err := conn.ReadJSON(&cmd); err != nil {
			return
		}

		switch strings.ToLower(strings.TrimSpace(cmd.Action)) {
		case "subscribe":
			if cmd.OrderID == uuid.Nil {
				continue
			}
			_ = a.Repo.UpsertSubscriber(r.Context(), cmd.OrderID, userID, connectionID)
		case "unsubscribe":
			if cmd.OrderID == uuid.Nil {
				continue
			}
			_ = a.Repo.DeactivateSubscriber(r.Context(), connectionID, &cmd.OrderID)
		}
	}
}
