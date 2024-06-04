package user

import (
	"errors"
	"fmt"
	"go-openapi/config"
	"go-openapi/pkg/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string `json:"email" gorm:"index"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
}

func (u *User) TableName() string {
	return "users"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&User{})
}

// ValidateUser 사용자 정보 유효성 검사
func ValidateUser(email, password string) (bool, error) {
	if email == "" || password == "" {
		return false, fmt.Errorf("invalid email or password")
	}
	if !validateEmail(email) {
		return false, fmt.Errorf("invalid email")
	}
	if !validatePassword(password) {
		return false, fmt.Errorf("invalid password(min: 6, max: 20, at least one letter and one number)")
	}
	return true, nil
}

// CreateUser 사용자 생성
func CreateUser(tx *gorm.DB, email, password string) (*User, error) {
	if ok, err := ValidateUser(email, password); !ok || err != nil {
		return nil, err
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	// 이미 존재하는 이메일인지 확인
	var existingUser User
	if err := tx.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	user := &User{
		Email:      email,
		Password:   hashedPassword,
		IsVerified: false,
	}
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
