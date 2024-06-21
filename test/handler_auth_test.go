package test

import (
	"encoding/base64"
	"go-openapi/cmd/api"
	"go-openapi/config"
	userModel "go-openapi/model/user"
	"go-openapi/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createAccessToken(t *testing.T, email string) string {
	// 사용자 생성
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/user/", nil)
	r.Form = map[string][]string{
		"email":    {email},
		"password": {"test1234"},
	}

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201, got %d", w.Code)

	// 인증
	db := config.GetDB()
	// 생성 여부 확인
	var instance userModel.User
	hashedEmail := utils.SHA256Email(email)
	if err := db.Where("hashed_email = ?", hashedEmail).First(&instance).Error; err != nil {
		t.Fail()
	}
	if instance.ID == 0 || instance.IsVerified {
		t.Fail()
	}
	instance.IsVerified = true
	db.Save(&instance)
	// 로그인
	r = httptest.NewRequest("POST", "/v1/auth/login/", nil)
	w = httptest.NewRecorder()
	r.Form = map[string][]string{
		"email":    {email},
		"password": {"test1234"},
	}

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	token := strings.Split(string(body), "\"accessToken\":\"")[1]
	token = strings.Split(token, "\"")[0]
	assert.NotEmpty(t, token)
	return token
}

// TestLoginHandlerSuccess 로그인 성공
func TestLoginHandlerSuccess(t *testing.T) {
	_ = createAccessToken(t, "test-auth-1@test.com")
}

func TestCreateClientHandlerSuccess(t *testing.T) {
	token := createAccessToken(t, "test-auth-2@test.com")

	// 클라이언트 생성
	r := httptest.NewRequest("POST", "/v1/auth/client/", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+token)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	// string(body) : {"clientId":"__CLIENT_ID__","clientSecret":"__CLIENT_SECRET__"}
	clientId := strings.Split(string(body), "\"clientId\":\"")[1]
	clientId = strings.Split(clientId, "\"")[0]
	clientSecret := strings.Split(string(body), "\"clientSecret\":\"")[1]
	clientSecret = strings.Split(clientSecret, "\"")[0]
	assert.NotEmpty(t, clientId)
	assert.NotEmpty(t, clientSecret)
}

// TestCreateTokenHandlerSuccess 토큰 생성 성공
func TestCreateTokenHandlerSuccess(t *testing.T) {
	token := createAccessToken(t, "test-auth-3@test.com")

	// 클라이언트 생성
	r := httptest.NewRequest("POST", "/v1/auth/client/", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+token)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	// string(body) : {"clientId":"__CLIENT_ID__","clientSecret":"__CLIENT_SECRET__"}
	clientId := strings.Split(string(body), "\"clientId\":\"")[1]
	clientId = strings.Split(clientId, "\"")[0]
	clientSecret := strings.Split(string(body), "\"clientSecret\":\"")[1]
	clientSecret = strings.Split(clientSecret, "\"")[0]
	assert.NotEmpty(t, clientId)
	assert.NotEmpty(t, clientSecret)

	r = httptest.NewRequest("POST", "/v1/auth/token/", nil)
	w = httptest.NewRecorder()
	req := `{"scope": "read:client", "grant_type": "client_credentials"}`
	r.Header.Set("Content-Type", "application/json")
	basicAuth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	r.Header.Set("Authorization", "Basic "+basicAuth)
	r.Body = ioutil.NopCloser(strings.NewReader(req))

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ = ioutil.ReadAll(w.Body)
	accessToken := strings.Split(string(body), "\"accessToken\":\"")[1]
	accessToken = strings.Split(accessToken, "\"")[0]
	assert.NotEmpty(t, accessToken)
}

// TestGetClientsHandlerSuccess 클라이언트 목록 조회 성공
func TestGetClientsHandlerSuccess(t *testing.T) {
	token := createAccessToken(t, "test-auth-4@test.com")

	// 클라이언트 생성
	r := httptest.NewRequest("POST", "/v1/auth/client/", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+token)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	// string(body) : {"clientId":"__CLIENT_ID__","clientSecret":"__CLIENT_SECRET__"}
	clientId := strings.Split(string(body), "\"clientId\":\"")[1]
	clientId = strings.Split(clientId, "\"")[0]
	clientSecret := strings.Split(string(body), "\"clientSecret\":\"")[1]
	clientSecret = strings.Split(clientSecret, "\"")[0]
	assert.NotEmpty(t, clientId)
	assert.NotEmpty(t, clientSecret)

	r = httptest.NewRequest("POST", "/v1/auth/token/", nil)
	w = httptest.NewRecorder()
	req := `{"scope": "read:client", "grant_type": "client_credentials"}`
	r.Header.Set("Content-Type", "application/json")
	basicAuth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	r.Header.Set("Authorization", "Basic "+basicAuth)
	r.Body = ioutil.NopCloser(strings.NewReader(req))

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ = ioutil.ReadAll(w.Body)
	accessToken := strings.Split(string(body), "\"accessToken\":\"")[1]
	accessToken = strings.Split(accessToken, "\"")[0]
	assert.NotEmpty(t, accessToken)

	// 클라이언트 목록 조회
	r = httptest.NewRequest("GET", "/v1/auth/client/", nil)
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+accessToken)

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ = ioutil.ReadAll(w.Body)
	assert.NotEmpty(t, string(body))
	// client 수 확인
	count := 0
	clients := strings.Split(string(body), "},{")
	for _, client := range clients {
		clientId := strings.Split(client, "\"clientId\":\"")[1]
		clientId = strings.Split(clientId, "\"")[0]
		clientSecret := strings.Split(client, "\"clientSecret\":\"")[1]
		clientSecret = strings.Split(clientSecret, "\"")[0]
		assert.NotEmpty(t, clientId)
		assert.NotEmpty(t, clientSecret)
		count++
	}
	if count == 0 {
		t.Fail()
	}
}

// TestDeleteClientsHandlerSuccess 클라이언트 삭제 성공
func TestDeleteClientsHandlerSuccess(t *testing.T) {
	token := createAccessToken(t, "test-auth-5@test.com")

	// 클라이언트 생성
	r := httptest.NewRequest("POST", "/v1/auth/client/", nil)
	w := httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+token)

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	// string(body) : {"clientId":"__CLIENT_ID__","clientSecret":"__CLIENT_SECRET__"}
	clientId := strings.Split(string(body), "\"clientId\":\"")[1]
	clientId = strings.Split(clientId, "\"")[0]
	clientSecret := strings.Split(string(body), "\"clientSecret\":\"")[1]
	clientSecret = strings.Split(clientSecret, "\"")[0]
	assert.NotEmpty(t, clientId)
	assert.NotEmpty(t, clientSecret)

	r = httptest.NewRequest("POST", "/v1/auth/token/", nil)
	w = httptest.NewRecorder()
	req := `{"scope": "read:client", "grant_type": "client_credentials"}`
	r.Header.Set("Content-Type", "application/json")
	basicAuth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	r.Header.Set("Authorization", "Basic "+basicAuth)
	r.Body = ioutil.NopCloser(strings.NewReader(req))

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ = ioutil.ReadAll(w.Body)
	accessToken := strings.Split(string(body), "\"accessToken\":\"")[1]
	accessToken = strings.Split(accessToken, "\"")[0]
	assert.NotEmpty(t, accessToken)

	// 클라이언트 목록 조회
	r = httptest.NewRequest("GET", "/v1/auth/client/", nil)
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+accessToken)

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)

	body, _ = ioutil.ReadAll(w.Body)
	assert.NotEmpty(t, string(body))
	// client 수 확인
	count := 0
	clients := strings.Split(string(body), "},{")
	lastId := 0
	for _, client := range clients {
		clientId := strings.Split(client, "\"clientId\":\"")[1]
		clientId = strings.Split(clientId, "\"")[0]
		clientSecret := strings.Split(client, "\"clientSecret\":\"")[1]
		clientSecret = strings.Split(clientSecret, "\"")[0]
		pk := strings.Split(client, "\"id\":\"")[1]
		pk = strings.Split(pk, "\"")[0]
		lastId, _ = strconv.Atoi(pk)
		assert.NotEmpty(t, clientId)
		assert.NotEmpty(t, clientSecret)
		count++
	}
	if count == 0 {
		t.Fail()
	}

	// 클라이언트 삭제
	r = httptest.NewRequest("DELETE", "/v1/auth/client/"+strconv.Itoa(lastId)+"/", nil)
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer "+accessToken)

	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)
}

// TestCreateTokenHandlerFail1 토큰 생성 실패(경로 오류)
func TestCreateTokenHandlerFail1(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/auth/tokens/", nil)
	w := httptest.NewRecorder()

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404, got %d", w.Code)
}

// TestCreateTokenHandlerFail2 토큰 생성 실패(요청 JSON 오류)
func TestCreateTokenHandlerFail2(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/auth/token/", nil)
	w := httptest.NewRecorder()
	req := `{"scopes": "read"}`
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(req))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"error":"invalid_request"}`, string(body))
}

// TestCreateTokenHandlerFail3 토큰 생성 실패(요청 JSON 오류)
func TestCreateTokenHandlerFail3(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/auth/token/", nil)
	w := httptest.NewRecorder()
	req := `{"scope": "read", "grant_type": "password"}`
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(req))

	server := api.Server()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400, got %d", w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"error":"invalid_request"}`, string(body))
}
