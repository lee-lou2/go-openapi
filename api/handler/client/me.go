package client

import (
	"github.com/gofiber/fiber/v3"
)

// GetMeHandler 내 정보 조회 핸들러
func GetMeHandler(c fiber.Ctx) error {
	clientId := fiber.Locals[uint](c, "client")
	return c.JSON(fiber.Map{
		"id": clientId,
	})
}
