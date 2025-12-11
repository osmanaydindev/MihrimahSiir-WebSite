package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupBooksRoutes(app *fiber.App) {

	app.Delete("/delete-book/:id", middlewares.IsAdmin, controllers.DeleteBook)
	app.Post("/create-book", middlewares.IsAdmin, controllers.CreateBook)
	app.Put("/update-book/:id", middlewares.IsAdmin, controllers.UpdateBook)
	app.Get("/get-books", controllers.GetBooks)
	app.Get("/get-books-paginated", controllers.GetBooksPaginated)
	app.Get("/get-book/:slug", controllers.GetBook)
	app.Get("/get-book-by-id/:id", controllers.GetBookById)
}
