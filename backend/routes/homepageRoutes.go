package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupHomepageRoutes(app *fiber.App) {
	// Public endpoint for users - filtered by role_id
	app.Get("/get-homepage-items", controllers.GetHomepageItems)

	// Admin endpoints - get all items without filtering
	app.Get("/get-all-homepage-items", middlewares.IsAdmin, controllers.GetAllHomepageItems)
	app.Get("/get-homepage-item/:id", middlewares.IsAdmin, controllers.GetHomepageItem)
	app.Post("/create-homepage-item", middlewares.IsAdmin, controllers.CreateHomepageItem)
	app.Put("/update-homepage-item/:id", middlewares.IsAdmin, controllers.UpdateHomepageItem)
	app.Delete("/delete-homepage-item/:id", middlewares.IsAdmin, controllers.DeleteHomepageItem)
}
