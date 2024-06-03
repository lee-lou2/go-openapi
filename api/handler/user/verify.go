package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/cmd/auth"
	authCmd "go-openapi/cmd/user"
	userCmd "go-openapi/cmd/user"
	"go-openapi/config"
	userModel "go-openapi/model/user"
)

// SendVerifyCodeHandler 인증 코드 전송 핸들러
func SendVerifyCodeHandler(c fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	if user.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already verified",
		})
	}
	// 인증 코드 생성
	if err := authCmd.SendVerifyCode(email); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Code sent"})
}

// VerifyCodeHandler 인증 코드 검증 핸들러
func VerifyCodeHandler(c fiber.Ctx) error {
	email := c.FormValue("email")
	code := c.FormValue("code")
	if code == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	// 코드 검증
	if !userCmd.VerifyCode(email, code) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid code",
		})
	}
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	user.IsVerified = true
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// 토큰 생성
	accessToken, refreshToken, err := auth.CreateTokenSet(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
