package main

import (
	"backend/database"
	"backend/middlewares"
	"backend/routes"
	"backend/services/mail"
	"backend/util"
	ws "backend/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// .env dosyasını yükle (yoksa ortam değişkenleri kullanılır)
	if err := godotenv.Load(); err != nil {
		println("Error loading .env file")
	}

	// JWT imzalama anahtarını yükle. Anahtar sabit olmazsa her restart'ta
	// tüm kullanıcı oturumları düşer; bu yüzden eksikse başlatma durur.
	if err := util.InitJWT(); err != nil {
		panic("JWT yapılandırması hatalı: " + err.Error())
	}

	// Upload klasörlerini hazırla (volume mount edilmiş olsa da boş olabilir)
	if err := os.MkdirAll("./uploads/profiles", 0o755); err != nil {
		panic("uploads klasörü oluşturulamadı: " + err.Error())
	}

	app := fiber.New(fiber.Config{
		// Limit request body size to prevent DoS attacks
		BodyLimit: 10 * 1024 * 1024, // 10MB
		// Add custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	database.ConnectDb()

	// Start WebSocket hub
	go ws.GlobalHub.Run()

	// Mail kuyruğu (tek worker; Resend'in 2 istek/sn limitini aşmamak için)
	mail.StartWorker()

	// Apply security headers middleware (FIRST - before any other middleware)
	app.Use(middlewares.SecurityHeaders())

	// Apply global rate limiter (100 requests per minute per IP)
	app.Use(middlewares.GlobalRateLimiter())

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     os.Getenv("CORS_ORIGINS"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cookie, Set-Cookie",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Setup routes
	routes.Setup(app)

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // default port
	}

	err := app.Listen(":" + port)
	if err != nil {
		panic("Could not listen to the server")
	}
}
