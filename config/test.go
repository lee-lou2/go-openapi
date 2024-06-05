package config

// IsTesting 테스트 환경인지 확인
func IsTesting() bool {
	return GetEnv("TESTING") == "true"
}
