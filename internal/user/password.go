package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	userPkg "go-openapi/pkg/user"
	"go-openapi/pkg/utils"
)

// VerifyCodeAndChangePassword 코드 검증 후 비밀번호 변경
func VerifyCodeAndChangePassword(email, code, password string) error {
	if !userPkg.VerifyCode(email, code, 2) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid code")
	}
	// 비밀번호 변경
	var user userModel.User
	db := config.GetDB()
	hashedEmail := utils.SHA256Email(email)
	if err := db.Where("hashed_email = ?", hashedEmail).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// SendPasswordVerifyCode 비밀번호 확인 코드 전송
func SendPasswordVerifyCode(email string) error {
	if err := userPkg.SendVerifyCode(email, 2); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
