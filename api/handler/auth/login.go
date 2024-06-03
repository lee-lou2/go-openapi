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
		errMsg := "Invalid request"
		if err != nil {
			errMsg = err.Error()
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": errMsg,
		})
	}
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	if err := db.Where("email = ?", email).Where("is_verified = true").First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	// 비밀번호 확인
	if !utils.CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}
	// 토큰 생성
	accessToken, refreshToken, err := authCmd.CreateTokenSet(user.ID)
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

// LogoutHandler 로그아웃 핸들러
func LogoutHandler(c fiber.Ctx) error {
	user := c.Locals("user").(uint)
	fmt.Println(user)
	return nil
}
