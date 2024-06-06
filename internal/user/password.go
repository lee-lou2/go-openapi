package user

import (
	"github.com/gofiber/fiber/v3"
	userModel "go-openapi/model/user"
	userPkg "go-openapi/pkg/user"
)

// VerifyCodeAndChangePassword 코드 검증 후 비밀번호 변경
func VerifyCodeAndChangePassword(email, code, password string) error {
	if !userPkg.VerifyCode(email, code, 2) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid code")
	}
	if err := userModel.UpdatePassword(email, password); err != nil {
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
