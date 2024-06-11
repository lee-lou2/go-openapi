package validation

// ValidateClientCredentials 클라이언트 자격 증명 유효성 검사
func ValidateClientCredentials(scope, grantType string) bool {
	if grantType == "" || scope == "" {
		return false
	}
	if grantType != "client_credentials" {
		return false
	}
	return true
}

// ValidateRefreshToken 토큰 갱신 유효성 검사
func ValidateRefreshToken(grantType, refreshToken string) bool {
	if grantType == "" || refreshToken == "" {
		return false
	}
	if grantType != "refresh_token" {
		return false
	}
	return true
}
