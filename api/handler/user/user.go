package user

import (
	"github.com/gofiber/fiber/v3"
	authCmd "go-openapi/cmd/user"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	"gorm.io/gorm"
)

// CreateUserHandler 사용자 생성 핸들러
func CreateUserHandler(c fiber.Ctx) error {
	body := new(struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	})
	if err := c.Bind().Body(body); err != nil {
		return err
	}

	// 트랜젝션 처리
	db := config.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": tx.Error.Error(),
		})
	}

	var user *userModel.User
	if err := tx.Transaction(func(tx *gorm.DB) error {
		// 사용자 생성
		var err error
		user, err = userModel.CreateUser(tx, body.Email, body.Password)
		if err != nil {
			return err
		}
		// 인증 코드 전송
		if err := authCmd.SendVerifyCode(user.Email, 1); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// 실패시 롤백
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	tx.Commit()

	return c.JSON(fiber.Map{
		"email": user.Email,
	})
}
