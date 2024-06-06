package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	userPkg "go-openapi/pkg/user"
	"gorm.io/gorm"
)

// CreateUser 사용자 생성
func CreateUser(email, password string) error {
	// 트랜젝션 처리
	db := config.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}

	var user *userModel.User
	if err := tx.Transaction(func(tx *gorm.DB) error {
		// 사용자 생성
		var err error
		user, err = userModel.CreateUser(tx, email, password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		// 인증 코드 전송
		if err := userPkg.SendVerifyCode(user.Email, 1); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return nil
	}); err != nil {
		// 실패시 롤백
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	tx.Commit()
	return nil
}
