package middleware

import (
	"github.com/gofiber/fiber/v3"
	authPkg "go-openapi/pkg/auth"
)

// AuthMiddleware 인증 미들웨어
func AuthMiddleware(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	user, err := authPkg.GetUser(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	c.Locals("user", user)
	return c.Next()
}
