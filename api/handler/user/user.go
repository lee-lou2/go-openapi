package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
)

// CreateUserHandler 사용자 생성 핸들러
func CreateUserHandler(c fiber.Ctx) error {
	body := new(struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	})
	if err := c.Bind().Body(body); err != nil {
		return err
	}
	if !validation.ValidateEmail(body.Email) || !validation.ValidatePassword(body.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	err := userInternal.CreateUser(body.Email, body.Password)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"email": body.Email,
	})
}
