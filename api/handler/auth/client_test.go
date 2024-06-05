package auth

import (
	clientModel "go-openapi/model/client"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"go-openapi/config"
)

func setupApp() *fiber.App {
	app := fiber.New()

	// 데이터베이스 초기화
	db := config.GetDB()
	_ = db.AutoMigrate(&clientModel.Client{})

	// 라우트 설정
	app.Post("/clients", CreateClientHandler, mockUserMiddleware)
	app.Get("/clients", GetClientsHandler, mockUserMiddleware)
	app.Delete("/clients/:id", DeleteClientHandler, mockUserMiddleware)

	return app
}

// mockUserMiddleware 사용자 인증 모킹 미들웨어
func mockUserMiddleware(c fiber.Ctx) error {
	// 실제 사용자 인증 로직으로 교체 필요
	fiber.Locals[uint](c, "user", 1)
	return c.Next()
}

func TestCreateClientHandler(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/clients", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetClientsHandler(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/clients", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestDeleteClientHandler(t *testing.T) {
	app := setupApp()

	// 우선 클라이언트를 생성하여 테스트 대상 클라이언트가 존재하도록 합니다.
	createReq := httptest.NewRequest("POST", "/clients", nil)
	createReq.Header.Set("Content-Type", "application/json")
	_, err := app.Test(createReq)
	assert.Nil(t, err)

	// 클라이언트 ID가 "1"이라고 가정하고 삭제 요청을 보냅니다.
	deleteReq := httptest.NewRequest("DELETE", "/clients/1", nil)
	deleteReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(deleteReq)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}
