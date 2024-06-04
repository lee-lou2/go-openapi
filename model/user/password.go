package user

import (
	"go-openapi/config"
	"go-openapi/pkg/utils"
)

// UpdatePassword 비밀번호 변경
func UpdatePassword(email, password string) error {
	// 비밀번호 변경
	var user User
	db := config.GetDB()
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return err
	}
	return nil
}
