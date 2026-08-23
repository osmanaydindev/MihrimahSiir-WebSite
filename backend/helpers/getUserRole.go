package helpers

import (
	"backend/database"
	"backend/models"
	"backend/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetUserIDFromCtx, token cookie'sinden kullanıcı id'sini çözer.
// Hata durumunda 0 döner. middlewares paketi controllers'ı import
// etmeden kullanıcı bazlı limit uygulayabilsin diye burada duruyor
// (helpers hiçbir üst katmanı import etmiyor, döngü riski yok).
func GetUserIDFromCtx(c *fiber.Ctx) uint {
	cookie := c.Cookies("token")
	if cookie == "" {
		return 0
	}
	userID, err := util.GetUserWithToken(cookie)
	if err != nil {
		return 0
	}
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

// GetUserRole retrieves the role ID of the authenticated user from the token
func GetUserRole(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("token")
	if cookie == "" {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "No token provided")
	}
	userID, err := util.GetUserWithToken(cookie)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	var admin models.Admin
	result := database.DB.Where("id = ?", userID).First(&admin)
	if result.Error != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}
	return admin.RoleID, nil
}
