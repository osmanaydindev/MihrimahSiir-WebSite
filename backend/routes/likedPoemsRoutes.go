package routes

import (
	"backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupLikedPoemsRoutes(app *fiber.App) {
	app.Post("/add-poem-to-liked/:id", controllers.AddLikedPoem)
	app.Post("/undo-poem-to-liked/:id", controllers.UndoLikedPoem)
	app.Get("/get-liked-poems-id/:id", controllers.GetLikedPoemsIdByAdminId)
	app.Get("/get-liked-poems-paginated/:id", controllers.GetLikedPoemsPaginated)
}
