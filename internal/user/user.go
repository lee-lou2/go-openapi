package user

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	userPkg "go-openapi/pkg/user"
	"go-openapi/pkg/utils"
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

	if err := tx.Transaction(func(tx *gorm.DB) error {
		// 사용자 생성
		var err error
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		// 이미 존재하는 이메일인지 확인
		var existingUser userModel.User
		hashedEmail := utils.SHA256Email(email)
		if err := tx.Where("hashed_email = ?", hashedEmail).First(&existingUser).Error; err == nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("email %s already exists", email))
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		aes := utils.NewAES256()
		encryptedEmail := aes.Encrypt(email)
		user := &userModel.User{
			HashedEmail:    hashedEmail,
			EncryptedEmail: encryptedEmail,
			Password:       hashedPassword,
			IsVerified:     false,
		}
		if err := tx.Create(user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		// 인증 코드 전송
		if err := userPkg.SendVerifyCode(email, 1); err != nil {
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
