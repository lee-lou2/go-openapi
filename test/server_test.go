package test

import (
	"go-openapi/cmd/api"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain 테스트 메인 설정
func TestMain(m *testing.M) {
	defer func() {
		// 테스트 데이터베이스 제거
		os.Remove("test.sqlite.db")
	}()
	m.Run()
}

// TestHealthCheck_Success 헬스체크 핸들러 테스트
func TestHealthCheck_Success(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)
}

// TestHealthCheck_Failed_Notfound 헬스체크 핸들러 테스트
func TestHealthCheck_Failed_Notfound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/_health", nil)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404, got %d", w.Code)
}
