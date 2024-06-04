package middleware

import (
	"github.com/gofiber/fiber/v3"
)

// PermissionMiddleware 권한 미들웨어
func PermissionMiddleware(scope string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		localScope := c.Locals("scopes")
		if localScope == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
		scopes := localScope.([]string)
		for _, s := range scopes {
			if s == scope {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}
