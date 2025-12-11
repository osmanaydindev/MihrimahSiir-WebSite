package routes

import (
	"backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupBooksReadRoutes(app *fiber.App) {
	app.Post("/add-book-to-reads/:id", controllers.AddBookToReads)
	app.Post("/delete-book-from-reads/:id", controllers.DeleteBookFromReads)
	app.Get("/get-reads-books-ids/:id", controllers.GetReadBooksIds)
	app.Get("/get-reads-books-paginated/:id", controllers.GetReadsBooksPaginated)
	app.Get("/get-not-reads-books-paginated/:id", controllers.GetNotReadsBooksPaginated)
}
