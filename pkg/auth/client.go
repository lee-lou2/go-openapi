package auth

import (
	"go-openapi/pkg/utils"
)

// GenerateClientKeys 클라이언트 키 생성
func GenerateClientKeys() (string, string, error) {
	clientID, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	clientSecret, err := utils.GenerateRandomString(64)
	if err != nil {
		return "", "", err
	}

	return clientID, clientSecret, nil
}
