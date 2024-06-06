package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
)

// SendVerifyCodeHandler 인증 코드 전송 핸들러
func SendVerifyCodeHandler(c fiber.Ctx) error {
	body := new(struct {
		Email string `json:"email"`
	})
	if err := c.Bind().JSON(body); err != nil {
		return err
	}
	if !validation.ValidateEmail(body.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}
	err := userInternal.ValidateUserAndSendVerifyCode(body.Email)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Code sent"})
}

// VerifyCodeHandler 인증 코드 검증 핸들러
func VerifyCodeHandler(c fiber.Ctx) error {
	code := fiber.Params[string](c, "code")
	body := new(struct {
		Email string `json:"email"`
	})
	if err := c.Bind().JSON(body); err != nil {
		return err
	}
	if !validation.ValidateEmail(body.Email) || !validation.ValidateCode(code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})

	}
	err := userInternal.VerifyCodeAndUpdateUser(body.Email, code)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "User verified"})
}
