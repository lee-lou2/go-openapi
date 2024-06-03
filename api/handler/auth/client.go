package auth

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/model/client"
)

// CreateClientKeysHandler 클라이언트 키 생성 핸들러
func CreateClientKeysHandler(c fiber.Ctx) error {
	user := c.Locals("user").(uint)
	instance, err := client.CreateClient(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"clientId":     instance.ClientId,
		"clientSecret": instance.ClientSecret,
	})
}
