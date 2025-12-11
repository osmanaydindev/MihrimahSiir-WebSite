package routes

import (
	"backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupCommentsRoutes(app *fiber.App) {
	app.Post("/add-comment", controllers.AddComment)
	app.Delete("/delete-comment/:comment_id", controllers.DeleteComment)
	//app.Post("/undo-bookmark/:id", controllers.UndoBookmark)
	//app.Get("/get-bookmark-id/:id", controllers.GetBookmarksIdByAdminId)
	//app.Get("/get-bookmark/:id", controllers.GetBookmarksByAdminId)
}
