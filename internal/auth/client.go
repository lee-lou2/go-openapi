package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	clientModel "go-openapi/model/client"
	"go-openapi/pkg/auth"
)

// CreateClient 클라이언트 생성
func CreateClient(userId uint) (*clientModel.Client, error) {
	if userId == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid user instance")
	}
	clientId, clientSecret, err := auth.GenerateClientKeys()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	db := config.GetDB()
	// 이미 존재하는 클라이언트 ID 또는 클라이언트 비밀키인지 확인
	var existingClient clientModel.Client
	if err := db.Where("client_id = ?", clientId).Or("client_secret = ?", clientSecret).First(&existingClient).Error; err == nil {
		return CreateClient(userId)
	}
	// 사용자별로 최대 10개까지 클라이언트 생성 가능
	var clientCount int64
	if err := db.Model(&clientModel.Client{}).Where("user_id = ?", userId).Count(&clientCount).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if clientCount >= 10 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "maximum client count exceeded")
	}
	scope := fmt.Sprintf("%s %s %s %s", clientModel.ScopeReadClient, clientModel.ScopeWriteClient, clientModel.ScopeReadToken, clientModel.ScopeReadDefault)
	client := &clientModel.Client{
		UserID:       userId,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        scope,
	}
	if err := db.Create(client).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return client, nil
}

// GetClients 클라이언트 조회
func GetClients(userId uint) (*[]clientModel.Client, error) {
	if userId == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid user instance")
	}
	db := config.GetDB()
	var clients []clientModel.Client
	if err := db.Where("user_id = ?", userId).Find(&clients).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return &clients, nil
}

// DeleteClient 클라이언트 삭제
func DeleteClient(userId uint, id string) error {
	if userId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user instance")
	}
	db := config.GetDB()
	var client clientModel.Client
	if err := db.Where("user_id = ?", userId).Where("id = ?", id).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "client not found")
	}
	if err := db.Delete(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
