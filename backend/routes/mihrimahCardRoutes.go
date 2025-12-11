package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupMihrimahCardRoutes(app *fiber.App) {
	// Public endpoint - all authenticated users can access
	app.Get("/get-mihrimah-cards", controllers.GetAllMihrimahCards)

	// Admin-only endpoints for management
	app.Get("/get-mihrimah-card/:id", middlewares.IsAdmin, controllers.GetMihrimahCard)
	app.Post("/create-mihrimah-card", middlewares.IsAdmin, controllers.CreateMihrimahCard)
	app.Put("/update-mihrimah-card/:id", middlewares.IsAdmin, controllers.UpdateMihrimahCard)
	app.Delete("/delete-mihrimah-card/:id", middlewares.IsAdmin, controllers.DeleteMihrimahCard)
}
