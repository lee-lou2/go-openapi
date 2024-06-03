package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"go-openapi/api/router"
	"go-openapi/config"
)

func Server() error {
	// Fiber 인스턴스 생성
	app := fiber.New()

	// 미들웨어
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.GetEnv("COOKIE_ENCRYPT_KEY"),
	}))

	// 라우터 설정
	router.BaseRouter(app)
	router.V1Router(app)

	// 미들웨어
	app.Use(func(c fiber.Ctx) error {
		// 404 처리
		return c.SendStatus(404)
	})

	// 서버 실행
	ServerPort := config.GetEnv("SERVER_PORT")
	return app.Listen(":" + ServerPort)
}
