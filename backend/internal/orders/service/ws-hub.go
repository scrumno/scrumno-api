package service

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type OrdersWebSocketHub interface {
	Register(connectionID string, conn *websocket.Conn)
	Unregister(connectionID string)
	Notify(connectionIDs []string, orderID uuid.UUID, status string)
}

type ordersWebSocketHub struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
}

func NewOrdersWebSocketHub() OrdersWebSocketHub {
	return &ordersWebSocketHub{
		connections: make(map[string]*websocket.Conn),
	}
}

func (h *ordersWebSocketHub) Register(connectionID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[connectionID] = conn
}

func (h *ordersWebSocketHub) Unregister(connectionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conn, ok := h.connections[connectionID]; ok {
		_ = conn.Close()
		delete(h.connections, connectionID)
	}
}

func (h *ordersWebSocketHub) Notify(connectionIDs []string, orderID uuid.UUID, status string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	payload, _ := json.Marshal(map[string]any{
		"order_id": orderID,
		"status":   status,
	})
	for _, connectionID := range connectionIDs {
		conn, ok := h.connections[connectionID]
		if !ok {
			continue
		}
		_ = conn.WriteMessage(websocket.TextMessage, payload)
	}
}
