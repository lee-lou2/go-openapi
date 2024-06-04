package auth

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/model/client"
)

// CreateClientHandler 클라이언트 키 생성 핸들러
func CreateClientHandler(c fiber.Ctx) error {
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

// GetClientsHandler 클라이언트 키 조회 핸들러
func GetClientsHandler(c fiber.Ctx) error {
	user := c.Locals("user").(uint)
	clients, err := client.GetClients(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	resp := make([]fiber.Map, 0)
	for _, instance := range *clients {
		resp = append(resp, fiber.Map{
			"id":           instance.ID,
			"clientId":     instance.ClientId,
			"clientSecret": instance.ClientSecret,
			"createdAt":    instance.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return c.JSON(resp)
}

// DeleteClientHandler 클라이언트 키 삭제 핸들러
func DeleteClientHandler(c fiber.Ctx) error {
	user := c.Locals("user").(uint)
	id := c.Params("id")
	err := client.DeleteClient(user, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}
