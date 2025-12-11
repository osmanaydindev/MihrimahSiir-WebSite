package routes

import (
	"backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthorRoutes(app *fiber.App) {
	// Public routes
	app.Get("/get-authors", controllers.GetAuthors)
	app.Get("/get-author/:slug", controllers.GetAuthor)
	app.Get("/get-all-authors-dropdown", controllers.GetAllAuthorsForDropdown)
	app.Get("/get-author-by-id/:id", controllers.GetAuthorById)

	// Admin routes (should be protected by middleware)
	app.Post("/create-author", controllers.CreateAuthor)
	app.Put("/update-author/:id", controllers.UpdateAuthor)
	app.Delete("/delete-author/:id", controllers.DeleteAuthor)
}
