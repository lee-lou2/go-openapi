package api

import (
	"go-openapi/api/middleware"
	"go-openapi/api/router"
	"net/http"
)

// func Server() error {
// 	// Fiber 인스턴스 생성
// 	engine := html.New("./views", ".html")
// 	app := fiber.New(fiber.Config{
// 		JSONEncoder: json.Marshal,
// 		JSONDecoder: json.Unmarshal,
// 		Views:       engine,
// 	})

// 	// 미들웨어 설정
// 	app.Use(requestid.New())
// 	// CORS
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "http://127.0.0.1, http://localhost",
// 		AllowHeaders: "Origin, Content-Type, Accept",
// 	}))
// 	// pprof
// 	app.Use(pprof.New())
// 	// 아이덤팅시 미들웨어
// 	app.Use(idempotency.New())
// 	// 쿠키 암호화
// 	app.Use(encryptcookie.New(encryptcookie.Config{
// 		Key: config.GetEnv("COOKIE_ENCRYPT_KEY"),
// 	}))
// 	// 로거 설정
// 	app.Use(logger.New(logger.Config{
// 		Format:     "${pid} ${status} - ${method} ${path}\n",
// 		TimeFormat: "2006-01-02 15:04:05",
// 		TimeZone:   "Asia/Seoul",
// 	}))

// 	// 라우터 설정
// 	router.BaseRouter(app)
// 	router.TemplateRouter(app)
// 	// router.V1Router(app)

// 	// 미들웨어
// 	app.Use(func(c fiber.Ctx) error {
// 		// 404 처리
// 		return c.SendStatus(404)
// 	})

// 	// 서버 실행
// 	ServerPort := config.GetEnv("SERVER_PORT")
// 	return app.Listen(":" + ServerPort)
// }

// Server 서버 설정
func Server() http.HandlerFunc {
	mux := http.NewServeMux()
	base := router.BaseRouter()
	mux.Handle("/", base)
	v1 := router.V1Router()
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	chain := middleware.Chain(
		middleware.LoggerMiddleware,
		middleware.RecoverMiddleware,
		middleware.CORSMiddleware,
		middleware.RequestIdMiddleware,
	)
	return chain(mux)
}
