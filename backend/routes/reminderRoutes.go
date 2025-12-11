package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ReminderRoutes(app *fiber.App) {
	// Public endpoints - all authenticated users can access based on permission
	app.Get("/reminders-paginated", controllers.GetRemindersPaginated)      // Hatırlatıcıları sayfalı getir (permission'a göre filtrelenmiş)
	app.Get("/reminders", controllers.GetAllReminders)                      // Tüm hatırlatıcıları getir (permission'a göre filtrelenmiş)
	app.Get("/reminders/:id", middlewares.IsAdmin, controllers.GetReminder) // Belirli bir hatırlatıcıyı getir

	// Admin-only endpoints for management
	app.Post("/reminders", middlewares.IsAdmin, controllers.CreateReminder)       // Yeni hatırlatıcı oluştur
	app.Put("/reminders/:id", middlewares.IsAdmin, controllers.UpdateReminder)    // Belirli bir hatırlatıcıyı güncelle
	app.Delete("/reminders/:id", middlewares.IsAdmin, controllers.DeleteReminder) // Belirli bir hatırlatıcıyı sil
}
