package routes

import (
	ws "backend/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupWebSocketRoutes(app *fiber.App) {
	app.Use("/ws", ws.WebSocketUpgrade)
	app.Get("/ws", websocket.New(ws.HandleWebSocket))
}
