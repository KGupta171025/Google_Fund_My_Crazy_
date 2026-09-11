package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"metro-synapse/backend/internal/broker"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local & cloud development
	},
}

// WSHandler manages WebSocket streaming connections
type WSHandler struct {
	broker *broker.EventBroker
}

// NewWSHandler returns a new WebSocket handler
func NewWSHandler(b *broker.EventBroker) *WSHandler {
	return &WSHandler{broker: b}
}

// HandleWS handles real-time WebSocket connection upgrades
func (ws *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	subID := fmt.Sprintf("ws-client-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
	sub := ws.broker.Subscribe(subID, "city:state")
	defer ws.broker.Unsubscribe(subID, "city:state")

	// Goroutine to drain incoming client pings/messages
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	// Stream events from broker to WebSocket client
	for event := range sub.Channel {
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := conn.WriteJSON(event.Payload); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}
