package helpers

import (
	"backend/database"
	"backend/models"
	"backend/util"
	"github.com/gofiber/fiber/v2"
)

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
