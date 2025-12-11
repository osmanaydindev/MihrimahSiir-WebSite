package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App) {

	app.Post("/create-admin", middlewares.AuthRateLimiter(), controllers.CreateAdmin)
	app.Post("/user", controllers.User)
	app.Post("/logout", controllers.LogOut)
	app.Get("/get-all-admins", controllers.GetAdmins)
	app.Get("/get-admins-management", controllers.GetAdminsForManagement)
	app.Get("/get-admin/:id", controllers.GetAdmin)
	app.Put("/update-admin/:id", controllers.UpdateAdmin)
	app.Delete("/delete-admin/:id", controllers.DeleteAdmin)

	// Profile routes with upload rate limiting (10 uploads per hour)
	app.Post("/upload-profile-image", middlewares.UploadRateLimiter(), controllers.UploadProfileImage)
	app.Get("/user-profile/:username", controllers.GetUserProfile)
	app.Put("/update-privacy", controllers.UpdatePrivacy)

	// Lazy loading routes for profile stats
	app.Get("/user-profile/:username/liked-poems", controllers.GetUserLikedPoems)
	app.Get("/user-profile/:username/read-books", controllers.GetUserReadBooks)
	app.Get("/user-profile/:username/bookmarked-poems", controllers.GetUserBookmarkedPoems)
	app.Get("/user-profile/:username/comments", controllers.GetUserComments)

}
