package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	authCmd "go-openapi/cmd/auth"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	"go-openapi/pkg/utils"
)

// LoginHandler 로그인 핸들러
func LoginHandler(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	if ok, err := userModel.ValidateUser(email, password); !ok || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed login",
		})
	}
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	if err := db.Where("email = ?", email).Where("is_verified = true").First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "failed login",
		})
	}
	// 비밀번호 확인
	if !utils.CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed login",
		})
	}
	// 토큰 생성(사용자 로그인의 경우 클라이언트 관리만 가능)
	accessToken, refreshToken, err := authCmd.CreateTokenSet(user.ID, "user", "read:client", "write:client")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed login",
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
