package config

import "os"

// GetEnv 환경 변수 조회
func GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv 환경 변수 설정
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}
