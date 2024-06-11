package client

import (
	"go-openapi/config"
	"go-openapi/model/user"
	"gorm.io/gorm"
)

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Client{})
}

type Client struct {
	gorm.Model
	User         user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	UserID       uint      `json:"user_id" gorm:"index"`
	ClientId     string    `json:"client_id" gorm:"index"`
	ClientSecret string    `json:"client_secret" gorm:"index"`
	Scope        string    `json:"scope"`
}

func (c *Client) TableName() string {
	return "clients"
}
