package user

import (
	"go-openapi/config"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	EncryptedEmail string `json:"encrypted_email" gorm:"index"`
	HashedEmail    string `json:"hashed_email" gorm:"index"`
	Password       string `json:"password"`
	Level          int    `json:"level"`
	IsVerified     bool   `json:"is_verified"`
}

func (u *User) TableName() string {
	return "users"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&User{})
}
