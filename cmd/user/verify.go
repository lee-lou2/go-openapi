package user

import (
	"go-openapi/config"
	"go-openapi/pkg/notify"
	"go-openapi/pkg/utils"
	"time"
)

// SendVerifyCode 인증 코드 전송
func SendVerifyCode(email string) error {
	code, err := utils.GenerateRandomString(4)
	if err != nil {
		return err
	}
	cache := config.GetCache()
	// RSA 암호화
	emailHash := utils.SHA256(email)
	cache.Set(emailHash, code, time.Duration(5*60))
	if err := notify.SendEmail(email, "인증 코드", code); err != nil {
		return err
	}
	return nil
}

// VerifyCode 인증 코드 확인
func VerifyCode(email, code string) bool {
	cache := config.GetCache()
	emailHash := utils.SHA256(email)
	cachedCode, ok := cache.Get(emailHash)
	if !ok {
		return false
	}
	// 캐시 삭제
	cache.Delete(emailHash)
	return cachedCode == code
}
