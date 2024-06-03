package config

import "os"

// GetEnv 환경 변수 조회
func GetEnv(key string) string {
	return os.Getenv(key)
}
