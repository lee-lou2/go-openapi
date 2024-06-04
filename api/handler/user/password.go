package user

import (
	"github.com/gofiber/fiber/v3"
	userCmd "go-openapi/cmd/user"
	userModel "go-openapi/model/user"
)

// SendPasswordResetCodeHandler 비밀번호 재설정 코드 전송 핸들러
func SendPasswordResetCodeHandler(c fiber.Ctx) error {
	body := new(struct {
		Email string `json:"email"`
	})
	if err := c.Bind().JSON(body); err != nil {
		return err
	}
	if err := userCmd.SendVerifyCode(body.Email, 2); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Code sent"})
}

// ResetPasswordHandler 비밀번호 재설정 핸들러
func ResetPasswordHandler(c fiber.Ctx) error {
	param := struct {
		Code string `uri:"code"`
	}{}
	err := c.Bind().URI(&param)
	if err != nil || param.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	body := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.Bind().JSON(body); err != nil {
		return err
	}
	if body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	if userCmd.VerifyCode(body.Email, param.Code, 2) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid code",
		})
	}
	if err := userModel.UpdatePassword(body.Email, body.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Password updated"})
}
