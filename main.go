package main

import (
	"github.com/joho/godotenv"
	"go-openapi/api"
	"log"
)

func main() {
	// 환경 변수 불러오기
	_ = godotenv.Load()
	// API 서버 실행
	log.Fatal(api.Server())
}
