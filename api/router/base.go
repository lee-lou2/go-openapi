package router

import "github.com/gofiber/fiber/v3"

// BaseRouter 기본 라우터
func BaseRouter(app fiber.Router) {
	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
}
