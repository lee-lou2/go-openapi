package auth

import (
	"github.com/gofiber/fiber/v3"
	authInternal "go-openapi/internal/auth"
)

// CreateClientHandler 클라이언트 키 생성 핸들러
func CreateClientHandler(c fiber.Ctx) error {
	user := fiber.Locals[uint](c, "user")
	instance, err := authInternal.CreateClient(user)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"clientId":     instance.ClientId,
		"clientSecret": instance.ClientSecret,
	})
}

// GetClientsHandler 클라이언트 키 조회 핸들러
func GetClientsHandler(c fiber.Ctx) error {
	user := fiber.Locals[uint](c, "user")
	clients, err := authInternal.GetClients(user)
	if err != nil {
		return err
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
	user := fiber.Locals[uint](c, "user")
	id := fiber.Params[string](c, "id")
	err := authInternal.DeleteClient(user, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}
