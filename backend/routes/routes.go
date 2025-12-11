package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// (5 attempts per 15 minutes)
	app.Post("/register", middlewares.AuthRateLimiter(), controllers.Register)
	app.Post("/login", middlewares.AuthRateLimiter(), controllers.Login)

	SetupWebSocketRoutes(app)

	app.Static("/uploads", "./uploads")

	app.Use(middlewares.IsAuthenticated)

	app.Get("/auth-check", controllers.AuthCheck)
	app.Get("/get-logs", middlewares.IsAdmin, controllers.GetLogs)
	app.Delete("/delete-log/:id", middlewares.IsAdmin, controllers.DeleteLog)

	SetupAdminRoutes(app)
	SetupLikedPoemsRoutes(app)
	SetupBookmarksRoutes(app)
	SetupPoemsRoutes(app)
	SetupBooksRoutes(app)
	SetupBooksReadRoutes(app)
	SetupCommentsRoutes(app)
	ReminderRoutes(app)
	SetupHomepageRoutes(app)
	SetupMihrimahCardRoutes(app)
	SetupFriendshipRoutes(app)
	SetupAuthorRoutes(app)

}

func SetupFriendshipRoutes(app *fiber.App) {
	// Send friend request
	app.Post("/send-friend-request", controllers.SendFriendRequest)

	// Get friend requests (received)
	app.Get("/get-friend-requests", controllers.GetFriendRequests)

	// Get sent requests
	app.Get("/get-sent-requests", controllers.GetSentRequests)

	// Get friends list
	app.Get("/get-friends", controllers.GetFriends)

	// Accept friend request
	app.Put("/accept-friend-request/:id", controllers.AcceptFriendRequest)

	// Reject friend request
	app.Delete("/reject-friend-request/:id", controllers.RejectFriendRequest)

	// Cancel sent friend request
	app.Delete("/cancel-friend-request/:id", controllers.CancelFriendRequest)

	// Remove friend
	app.Delete("/remove-friend/:id", controllers.RemoveFriend)

	// Get pending requests count (for badge)
	app.Get("/get-pending-requests-count", controllers.GetPendingRequestsCount)
}
