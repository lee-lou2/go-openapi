package client

import (
	"errors"
	"go-openapi/config"
	"go-openapi/model/user"
	"go-openapi/pkg/auth"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	User         user.User `gorm:"foreignKey:UserID"`
	UserID       uint      `json:"user_id"`
	ClientId     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
}

func (c *Client) TableName() string {
	return "clients"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Client{})
}

// CreateClient 클라이언트 생성
func CreateClient(userId uint) (*Client, error) {
	if userId == 0 {
		return nil, errors.New("invalid user instance")
	}
	clientId, clientSecret, err := auth.GenerateClientKeys()
	if err != nil {
		return nil, err
	}
	db := config.GetDB()
	// 이미 존재하는 클라이언트 ID 또는 클라이언트 비밀키인지 확인
	var existingClient Client
	if err := db.Where("client_id = ?", clientId).Or("client_secret = ?", clientSecret).First(&existingClient).Error; err == nil {
		return CreateClient(userId)
	}
	// 사용자별로 최대 10개까지 클라이언트 생성 가능
	var clientCount int64
	if err := db.Model(&Client{}).Where("user_id = ?", userId).Count(&clientCount).Error; err != nil {
		return nil, err
	}
	if clientCount >= 10 {
		return nil, errors.New("maximum number of clients reached")
	}
	client := &Client{
		UserID:       userId,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
	if err := db.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

// GetClients 클라이언트 조회
func GetClients(userId uint) (*[]Client, error) {
	if userId == 0 {
		return nil, errors.New("invalid user instance")
	}
	db := config.GetDB()
	var clients []Client
	if err := db.Where("user_id = ?", userId).Find(&clients).Error; err != nil {
		return nil, err
	}
	if len(clients) == 0 {
		return nil, errors.New("client not found")
	}
	return &clients, nil
}

// DeleteClient 클라이언트 삭제
func DeleteClient(userId uint, id string) error {
	if userId == 0 {
		return errors.New("invalid user instance")
	}
	db := config.GetDB()
	var client Client
	if err := db.Where("user_id = ?", userId).Where("id = ?", id).First(&client).Error; err != nil {
		return err
	}
	if err := db.Delete(&client).Error; err != nil {
		return err
	}
	return nil
}
