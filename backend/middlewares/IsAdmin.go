package middlewares

import (
	"backend/database"
	"backend/models"
	"backend/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func IsAdmin(c *fiber.Ctx) error {
	// Get token from cookie
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Extract user ID from token
	userID, err := util.GetUserWithToken(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
			"error":   err,
		})
	}

	// Convert userID string to uint
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid user ID",
		})
	}

	// Get user from database
	var admin models.Admin
	result := database.DB.Where("id = ?", id).First(&admin)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// Check if user is admin (role_id = 1)
	if admin.RoleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized: admin access required",
		})
	}

	// User is admin, proceed
	return c.Next()
}
