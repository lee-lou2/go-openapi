package api

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"go-openapi/api/router"
	"go-openapi/config"
)

func Server() error {
	// Fiber 인스턴스 생성
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Views:       engine,
	})

	// 미들웨어 설정
	app.Use(requestid.New())
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1, http://localhost",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// recover
	app.Use(recover.New())
	// pprof
	app.Use(pprof.New())
	// 아이덤팅시 미들웨어
	app.Use(idempotency.New())
	// 쿠키 암호화
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.GetEnv("COOKIE_ENCRYPT_KEY"),
	}))
	// 로거 설정
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Seoul",
	}))

	// 라우터 설정
	router.BaseRouter(app)
	router.TemplateRouter(app)
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
