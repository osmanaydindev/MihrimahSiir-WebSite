package websocket

import (
	"backend/util"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"strconv"
)

// WebSocketUpgrade middleware to check if connection is websocket upgrade
func WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *websocket.Conn) {
	// Get user ID from query parameter (token)
	token := c.Query("token")
	if token == "" {
		log.Println("WebSocket: No token provided")
		c.Close()
		return
	}

	userIDStr, err := util.GetUserWithToken(token)
	if err != nil {
		log.Println("WebSocket: Invalid token")
		c.Close()
		return
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("WebSocket: Invalid user ID")
		c.Close()
		return
	}

	userID := uint(userIDInt)

	// Create client
	client := &Client{
		UserID: userID,
		Conn:   make(chan []byte, 256),
	}

	// Register client
	GlobalHub.Register <- client

	// Cleanup on disconnect
	defer func() {
		GlobalHub.Unregister <- client
	}()

	// Start goroutine to write messages to websocket
	go func() {
		for {
			message, ok := <-client.Conn
			if !ok {
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		}
	}()

	// Read messages from client (handle ping/pong and keep connection alive)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		// Parse message to check if it's a ping
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
				// Respond with pong
				pongMsg := []byte(`{"type":"pong"}`)
				if err := c.WriteMessage(websocket.TextMessage, pongMsg); err != nil {
					log.Println("WebSocket pong error:", err)
					break
				}
			}
		}
	}
}
