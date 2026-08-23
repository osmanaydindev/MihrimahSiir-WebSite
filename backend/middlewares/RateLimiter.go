package middlewares

import (
	"backend/helpers"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Max        int
	Expiration time.Duration
}

// UserRateLimiter, IP yerine kullanıcı id'sine göre sayar. Diğer
// limiterlar sadece c.IP() kullanıyor; aynı NAT arkasındaki kullanıcılar
// birbirinin kotasını tüketiyordu. Kimlik çözülemezse IP'ye düşer.
//
// Not: store bellekte tutuluyor, yani restart'ta sıfırlanır. Asıl
// koruma controller'daki DB tabanlı sayımlarda.
func UserRateLimiter(config RateLimiterConfig, message string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		// Sadece başarılı istekler sayılır. Asıl maliyet (dış API çağrısı +
		// mail) yalnızca 2xx dönen isteklerde oluşuyor; hatalı ISBN yazan
		// kullanıcının bir saat kilitlenmesi gereksiz. Ucuz hatalı istekler
		// zaten global IP limitiyle (100/dk) sınırlı.
		SkipFailedRequests: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := helpers.GetUserIDFromCtx(c); userID > 0 {
				return "u:" + strconv.Itoa(int(userID))
			}
			return "ip:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": message,
			})
		},
	})
}

// GlobalRateLimiter creates a global rate limiter
func GlobalRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // 100 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// AuthRateLimiter creates a rate limiter for authentication endpoints
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,                // 5 attempts
		Expiration: 15 * time.Minute, // per 15 minutes
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many authentication attempts. Please try again in 15 minutes.",
			})
		},
	})
}

// UploadRateLimiter creates a rate limiter for file upload endpoints
func UploadRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,            // 10 uploads
		Expiration: 1 * time.Hour, // per hour
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Upload limit exceeded. Please try again later.",
			})
		},
	})
}

// CustomRateLimiter creates a custom rate limiter with specified config
func CustomRateLimiter(config RateLimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}
