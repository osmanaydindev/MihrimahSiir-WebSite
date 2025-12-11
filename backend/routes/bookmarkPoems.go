package routes

import (
	"backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupBookmarksRoutes(app *fiber.App) {
	app.Post("/add-bookmark/:id", controllers.AddBookmark)
	app.Post("/undo-bookmark/:id", controllers.UndoBookmark)
	app.Get("/get-bookmark-id/:id", controllers.GetBookmarksIdByAdminId)
	app.Get("/get-bookmarks-paginated/:id", controllers.GetBookmarksPaginated)
}
