package notify

import (
	"errors"
	"go-openapi/config"
	"gopkg.in/mail.v2"
	"strconv"
)

// SendEmail 이메일 발송
func SendEmail(to, subject, body string) error {
	// SMTP 서버 설정
	smtpHost := config.GetEnv("EMAIL_SMTP_HOST")
	smtpPortString := config.GetEnv("EMAIL_SMTP_PORT")
	username := config.GetEnv("EMAIL_USERNAME")
	password := config.GetEnv("EMAIL_PASSWORD")
	if smtpHost == "" || smtpPortString == "" || username == "" || password == "" {
		return errors.New("SMTP server configuration is missing")
	}
	smtpPort, err := strconv.Atoi(smtpPortString)
	if err != nil {
		return err
	}

	// 메일 내용 설정
	from := username

	// 메시지 생성
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 다이얼러 설정
	d := mail.NewDialer(smtpHost, smtpPort, username, password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// 이메일 발송
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
