package middlewares

import (
	"backend/util"
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	_, err := util.GetUserWithToken(cookie)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
			"error":   err,
		})
	}
	return c.Next()
}

//func IsAuthenticated(c *fiber.Ctx) error {
//	token := c.Get("Authorization")
//	if token == "" {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"message": "missing token",
//		})
//	}
//
//	// "Bearer " prefixini kaldır
//	if len(token) > 7 && token[:7] == "Bearer " {
//		token = token[7:]
//	}
//
//	_, err := util.GetUserWithToken(token)
//	if err != nil {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"message": "unauthenticated",
//			"error":   err.Error(),
//		})
//	}
//
//	return c.Next()
//}
