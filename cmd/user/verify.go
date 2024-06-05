package user

import (
	"encoding/base64"
	"fmt"
	"go-openapi/config"
	"go-openapi/pkg/notify"
	"go-openapi/pkg/utils"
	"strconv"
	"strings"
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
	// 이메일 주소와 코드를 결합하여 인코딩
	combined := email + ":" + code
	encoded := base64.URLEncoding.EncodeToString([]byte(combined))
	encoded = strings.TrimRight(encoded, "=")
	host := config.GetEnv("SERVER_HOST")
	port := config.GetEnv("SERVER_PORT")
	if err := notify.SendEmail(email, subject, fmt.Sprintf("%s:%s/verify/%d?code=%s", host, port, codeType, encoded)); err != nil {
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
