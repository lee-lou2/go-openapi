package validation

// ValidateCode 코드 유효성 검사
func ValidateCode(code string) bool {
	if code == "" {
		return false
	}
	if len(code) != 8 {
		return false
	}
	return true
}
