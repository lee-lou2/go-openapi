package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"go-openapi/config"
	"strings"
)

// SHA256 SHA-256 해시 생성
func SHA256(data string) string {
	dataWithSalt := data + config.GetEnv("SHA256_SALT")
	hash := sha256.New()
	hash.Write([]byte(dataWithSalt))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}

// SHA256Email 이메일 주소를 SHA-256 해시로 변환
func SHA256Email(email string) string {
	// @ 기준으로 앞/뒤 각각 해싱하여 합침
	emailParts := strings.Split(email, "@")
	return SHA256(emailParts[0]) + ":" + SHA256(emailParts[1])
}
