package test

import (
	"go-openapi/cmd/api"
	"go-openapi/config"
	"go-openapi/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userModel "go-openapi/model/user"

	"github.com/stretchr/testify/assert"
)

// TestCreateUserSuccess 사용자 생성 성공
func TestCreateUserSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Form = map[string][]string{
		"email":    {"test@test.com"},
		"password": {"test1234"},
	}

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"email":"test@test.com"}`, string(body))

	db := config.GetDB()
	// 생성 여부 확인
	var instance userModel.User
	hashedEmail := utils.SHA256Email("test@test.com")
	if err := db.Where("hashed_email = ?", hashedEmail).First(&instance).Error; err != nil {
		t.Fail()
	}
	if instance.ID == 0 || instance.IsVerified {
		t.Fail()
	}
}

// TestCreateUserFail1 사용자 생성 실패(경로 오류)
func TestCreateUserFail(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/users", nil)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404, got %d", w.Code)
}

// TestCreateUserFail2 사용자 생성 실패(요청 JSON 오류)
func TestCreateUserFail2(t *testing.T) {
	w := httptest.NewRecorder()
	user := `{"email": "test", "password": "test"}`
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Invalid request"}`, string(body))
}

// TestCreateUserFail3 사용자 생성 실패(요청 JSON 오류)
func TestCreateUserFail3(t *testing.T) {
	w := httptest.NewRecorder()
	user := `{"email": "test", "password": "test123"}`
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Invalid request"}`, string(body))
}

// TestCreateUserFail4 사용자 생성 실패(요청 JSON 오류)
func TestCreateUserFail4(t *testing.T) {
	w := httptest.NewRecorder()
	user := `{"email": "test@test.com", "password": "test"}`
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Invalid request"}`, string(body))
}

// TestCreateUserFail5 사용자 생성 실패(이메일 형식 오류)
func TestCreateUserFail5(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Form = map[string][]string{
		"email":    {"test"},
		"password": {"test1234"},
	}

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Invalid request"}`, string(body))
}

// TestCreateUserFail6 사용자 생성 실패(비밀번호 형식 오류)
func TestCreateUserFail6(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Form = map[string][]string{
		"email":    {"test@test.com"},
		"password": {"123"},
	}

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Invalid request"}`, string(body))
}

// TestSendVerifyCodeHandlerSuccess 인증 코드 전송 성공
func TestSendVerifyCodeHandlerSuccess(t *testing.T) {
	db := config.GetDB()
	// 생성 여부 확인
	var instance userModel.User
	hashedEmail := utils.SHA256Email("test@test.com")
	if err := db.Where("hashed_email = ?", hashedEmail).First(&instance).Error; err != nil {
		t.Fail()
	}
	if instance.ID == 0 || instance.IsVerified {
		t.Fail()
	}

	w := httptest.NewRecorder()
	user := `{"email": "test@test.com"}`
	r := httptest.NewRequest("POST", "/v1/user/verify", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Code sent"}`, string(body))
}

// TestSendVerifyCodeHandlerFail1 인증 코드 전송 실패(사용자 없음)
func TestSendVerifyCodeHandlerFail1(t *testing.T) {
	w := httptest.NewRecorder()
	user := `{"email": "test2@test.com"}`
	r := httptest.NewRequest("POST", "/v1/user/verify", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"User not found"}`, string(body))
}

// TestSendVerifyCodeHandlerFail2 인증 코드 전송 실패(이미 인증된 사용자)
func TestSendVerifyCodeHandlerFail2(t *testing.T) {
	db := config.GetDB()
	// 생성 여부 확인
	var instance userModel.User
	hashedEmail := utils.SHA256Email("test@test.com")
	if err := db.Where("hashed_email = ?", hashedEmail).First(&instance).Error; err != nil {
		t.Fail()
	}
	if instance.ID == 0 || instance.IsVerified {
		t.Fail()
	}
	instance.IsVerified = true
	db.Save(&instance)
	w := httptest.NewRecorder()
	user := `{"email": "test@test.com"}`
	r := httptest.NewRequest("POST", "/v1/user/verify", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"User already verified"}`, string(body))
}

// TestVerifyCodeHandlerSuccess 인증 코드 검증 성공
func TestVerifyCodeHandlerSuccess(t *testing.T) {

}

func TestVerifyCodeHandlerFail1(t *testing.T) {

}

// TestSendPasswordResetCodeHandlerSuccess 비밀번호 재설정 코드 전송 성공
func TestSendPasswordResetCodeHandlerSuccess(t *testing.T) {
	db := config.GetDB()
	// 생성 여부 확인
	var instance userModel.User
	hashedEmail := utils.SHA256Email("test@test.com")
	if err := db.Where("hashed_email = ?", hashedEmail).First(&instance).Error; err != nil {
		t.Fail()
	}
	instance.IsVerified = true
	db.Save(&instance)
	if instance.ID == 0 || !instance.IsVerified {
		t.Fail()
	}

	w := httptest.NewRecorder()
	user := `{"email": "test@test.com"}`
	r := httptest.NewRequest("POST", "/v1/user/password", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"Code sent"}`, string(body))
}

func TestSendPasswordResetCodeHandlerFail1(t *testing.T) {
	w := httptest.NewRecorder()
	user := `{"email": "test2@test.com"}`
	r := httptest.NewRequest("POST", "/v1/user/password", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"User not found"}`, string(body))
}

func TestSendPasswordResetCodeHandlerFail2(t *testing.T) {
	db := config.GetDB()
	// 생성 여부 확인
	var instance userModel.User
	hashedEmail := utils.SHA256Email("test@test.com")
	if err := db.Where("hashed_email = ?", hashedEmail).First(&instance).Error; err != nil {
		t.Fail()
	}
	if instance.ID == 0 || !instance.IsVerified {
		t.Fail()
	}
	instance.IsVerified = false
	db.Save(&instance)
	w := httptest.NewRecorder()
	user := `{"email": "test@test.com"}`
	r := httptest.NewRequest("POST", "/v1/user/password", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(user))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"User not verified"}`, string(body))
}

// TestResetPasswordHandlerSuccess 비밀번호 재설정 성공
func TestResetPasswordHandlerSuccess(t *testing.T) {

}

func TestResetPasswordHandlerFail1(t *testing.T) {

}
