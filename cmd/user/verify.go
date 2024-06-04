package user

import (
	"go-openapi/config"
	"go-openapi/pkg/notify"
	"go-openapi/pkg/utils"
	"strconv"
	"time"
)

// SendVerifyCode 인증 코드 전송
func SendVerifyCode(email string, codeType int) error {
	code, err := utils.GenerateRandomString(8)
	if err != nil {
		return err
	}
	cache := config.GetCache()
	// RSA 암호화
	emailHash := utils.SHA256(email)
	cache.Set(emailHash+":"+strconv.Itoa(codeType), code, time.Duration(5*60))
	subject := ""
	switch codeType {
	case 1:
		subject = "인증 코드"
	case 2:
		subject = "비밀번호 재설정 코드"
	}
	if err := notify.SendEmail(email, subject, code); err != nil {
		return err
	}
	return nil
}

// VerifyCode 인증 코드 확인
func VerifyCode(email, code string, codeType int) bool {
	cache := config.GetCache()
	emailHash := utils.SHA256(email)
	cachedCode, ok := cache.Get(emailHash + ":" + strconv.Itoa(codeType))
	if !ok {
		return false
	}
	// 캐시 삭제
	cache.Delete(emailHash + ":" + strconv.Itoa(codeType))
	return cachedCode == code
}
