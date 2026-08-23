package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupBookRequestRoutes(app *fiber.App) {
	// Kullanıcı tarafı. Talep oluşturma her seferinde bir dış API
	// çağrısı ve bir mail tetiklediği için kullanıcı bazlı limitli.
	app.Post("/create-book-request",
		middlewares.UserRateLimiter(
			middlewares.RateLimiterConfig{Max: 5, Expiration: 1 * time.Hour},
			"Çok fazla kitap isteği gönderdiniz. Lütfen bir saat sonra tekrar deneyin.",
		),
		controllers.CreateBookRequest)
	app.Get("/get-my-book-requests", controllers.GetMyBookRequests)
	app.Delete("/cancel-book-request/:id", controllers.CancelBookRequest)

	// Admin tarafı
	app.Get("/get-book-requests", middlewares.IsAdmin, controllers.GetBookRequests)
	app.Get("/get-book-request-count", middlewares.IsAdmin, controllers.GetBookRequestCount)
	app.Post("/approve-book-request/:id", middlewares.IsAdmin, controllers.ApproveBookRequest)
	app.Post("/reject-book-request/:id", middlewares.IsAdmin, controllers.RejectBookRequest)
	app.Post("/refresh-book-request/:id", middlewares.IsAdmin, controllers.RefreshBookRequest)
}
