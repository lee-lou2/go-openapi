package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/validation"
	authInternal "go-openapi/internal/auth"
	"log"
)

// LoginHandler 로그인 핸들러
func LoginHandler(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	if !validation.ValidateEmail(email) || !validation.ValidatePassword(password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	accessToken, refreshToken, err := authInternal.GetTokenFromLogin(email, password)
	if err != nil {
		// 모든 오류 내용 통일
		log.Println(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to login",
		})
	}
	return c.JSON(fiber.Map{
		"tokenType":             "Bearer",
		"accessToken":           accessToken,
		"refreshToken":          refreshToken,
		"accessTokenExpiresIn":  3600,
		"refreshTokenExpiresIn": 86400,
	})
}

// LogoutHandler 로그아웃 핸들러
func LogoutHandler(c fiber.Ctx) error {
	user := fiber.Locals[uint](c, "user")
	fmt.Println(user)
	return nil
}
