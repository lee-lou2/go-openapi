package validation

import (
	"regexp"
	"unicode"
)

// ValidateEmail 이메일 유효성 검사
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}

// ValidatePassword 비밀번호 유효성 검사
func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	var hasLetter, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
