package routes

import (
	"time"

	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupFeedRoutes, routes.Setup içinde IsAuthenticated'ten SONRA çağrılır —
// akışın tamamı oturum gerektirir.
func SetupFeedRoutes(app *fiber.App) {
	// Beğen/kaydet uçları kullanıcı başına sınırlanıyor: tek istek ucuz ama
	// döngüye alınırsa gereksiz yazma yükü üretir.
	actionLimiter := middlewares.UserRateLimiter(middlewares.RateLimiterConfig{
		Max:        60,
		Expiration: time.Minute,
	}, "Çok fazla işlem yaptınız, lütfen biraz bekleyin.")

	app.Get("/feed", controllers.GetFeed)
	app.Get("/saved-comments", controllers.GetSavedComments)

	app.Post("/comments/:comment_id/like", actionLimiter, controllers.LikeComment)
	app.Delete("/comments/:comment_id/like", actionLimiter, controllers.UnlikeComment)
	app.Post("/comments/:comment_id/save", actionLimiter, controllers.SaveComment)
	app.Delete("/comments/:comment_id/save", actionLimiter, controllers.UnsaveComment)
}
