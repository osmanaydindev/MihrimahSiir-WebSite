package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupPoemsRoutes(app *fiber.App) {
	app.Post("/create-poem", middlewares.IsAdmin, controllers.CreatePoem)
	app.Delete("/delete-poem/:id", middlewares.IsAdmin, controllers.DeletePoem)
	app.Put("/update-poem/:id", middlewares.IsAdmin, controllers.UpdatePoem)
	app.Get("/get-poems", controllers.GetPoemsPaginated)
	app.Get("/get-poem/:slug", controllers.GetPoem)
	app.Get("/get-poem-by-id/:id", controllers.GetPoemById)
	app.Get("/get-all-poems", controllers.GetAllPoems)
	app.Get("/get-latest-poems", controllers.GetLatestPoems)
	app.Get("/get-search-poems", controllers.GetSearchPoems)
	app.Get("/get-popular-poems", controllers.GetPopularPoems)
}
