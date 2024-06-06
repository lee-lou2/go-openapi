package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
)

// SendPasswordResetCodeHandler 비밀번호 재설정 코드 전송 핸들러
func SendPasswordResetCodeHandler(c fiber.Ctx) error {
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
	err := userInternal.SendPasswordVerifyCode(body.Email)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Code sent"})
}

// ResetPasswordHandler 비밀번호 재설정 핸들러
func ResetPasswordHandler(c fiber.Ctx) error {
	body := new(struct {
		Code     string `uri:"code"`
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.Bind().JSON(body); err != nil {
		return err
	}
	if !validation.ValidateEmail(body.Email) || !validation.ValidatePassword(body.Password) || !validation.ValidateCode(body.Code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	err := userInternal.VerifyCodeAndChangePassword(body.Email, body.Code, body.Password)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Password updated"})
}
