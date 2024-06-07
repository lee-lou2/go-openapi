package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	userPkg "go-openapi/pkg/user"
	"go-openapi/pkg/utils"
)

// ValidateUserAndSendVerifyCode 사용자 확인 및 인증 코드 전송
func ValidateUserAndSendVerifyCode(email string) error {
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	hashedEmail := utils.SHA256Email(email)
	if err := db.Where("hashed_email = ?", hashedEmail).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User not found")
	}
	if user.IsVerified {
		return fiber.NewError(fiber.StatusBadRequest, "User already verified")
	}
	// 인증 코드 생성
	if err := userPkg.SendVerifyCode(email, 1); err != nil {
		return err
	}
	return nil
}

// VerifyCodeAndUpdateUser 인증 코드 확인 및 사용자 업데이트
func VerifyCodeAndUpdateUser(email, code string) error {
	// 코드 검증
	if !userPkg.VerifyCode(email, code, 1) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid code")
	}
	// 사용자 조회
	db := config.GetDB()
	user := userModel.User{}
	hashedEmail := utils.SHA256Email(email)
	if err := db.Where("hashed_email = ?", hashedEmail).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User not found")
	}
	user.IsVerified = true
	if err := db.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
