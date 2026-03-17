package handlers

import (
	"geofence/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for real-time alerts
type WebSocketHandler struct {
	alertPublisher domain.AlertPublisher
	upgrader       websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(alertPublisher domain.AlertPublisher) *WebSocketHandler {
	return &WebSocketHandler{
		alertPublisher: alertPublisher,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now, restrict in production
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// HandleWebSocket handles WebSocket connections at /ws/alerts
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Create alert channel
	alertChan := make(chan *domain.RealTimeAlert, 100)

	// Subscribe to alerts
	subscriptionID := h.alertPublisher.Subscribe(alertChan)
	defer h.alertPublisher.Unsubscribe(subscriptionID)

	// Handle incoming messages (for heartbeat/ping)
	go func() {
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
		}
	}()

	// Send alerts to client
	for alert := range alertChan {
		if err := conn.WriteJSON(alert); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}
