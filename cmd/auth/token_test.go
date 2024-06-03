package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go-openapi/config"
	"testing"
	"time"
)

// 환경 설정 초기화
func init() {
	config.SetEnv("JWT_SECRET", "your_jwt_secret")
}

// TestCreateTokenSet는 CreateTokenSet 함수의 동작을 테스트합니다.
func TestCreateTokenSet(t *testing.T) {
	userId := uint(1)

	accessToken, refreshToken, err := CreateTokenSet(userId)
	assert.NoError(t, err)           // 오류가 없는지 확인
	assert.NotEmpty(t, accessToken)  // 액세스 토큰이 비어있지 않은지 확인
	assert.NotEmpty(t, refreshToken) // 리프레시 토큰이 비어있지 않은지 확인

	// 액세스 토큰 검증
	accessClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "access", accessClaims["token_type"])  // 토큰 타입이 "access"인지 확인
	assert.Equal(t, float64(userId), accessClaims["user"]) // 사용자 ID가 일치하는지 확인

	// 리프레시 토큰 검증
	refreshClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "refresh", refreshClaims["token_type"]) // 토큰 타입이 "refresh"인지 확인
	assert.Equal(t, float64(userId), refreshClaims["user"]) // 사용자 ID가 일치하는지 확인
}

// TestCreateToken 은 CreateToken 함수의 동작을 다양한 케이스로 테스트합니다.
func TestCreateToken(t *testing.T) {
	tests := []struct {
		tokenType string
		userId    uint
		exp       int
	}{
		{"test", 1, 3600},
		{"access", 2, 7200},
		{"refresh", 3, 86400},
	}

	for _, tt := range tests {
		t.Run(tt.tokenType, func(t *testing.T) {
			token, err := CreateToken(tt.tokenType, tt.userId, tt.exp)
			assert.NoError(t, err)    // 오류가 없는지 확인
			assert.NotEmpty(t, token) // 토큰이 비어있지 않은지 확인

			claims := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.GetEnv("JWT_SECRET")), nil
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.tokenType, claims["token_type"])                                                                                    // 토큰 타입이 일치하는지 확인
			assert.Equal(t, float64(tt.userId), claims["user"])                                                                                    // 사용자 ID가 일치하는지 확인
			assert.WithinDuration(t, time.Now().Add(time.Second*time.Duration(tt.exp)), time.Unix(int64(claims["exp"].(float64)), 0), time.Second) // 만료 시간이 정확한지 확인
		})
	}
}

// TestCreateToken_InvalidSecret 은 잘못된 비밀키를 사용하여 CreateToken 함수가 오류를 반환하는지 테스트합니다.
func TestCreateToken_InvalidSecret(t *testing.T) {
	tokenType := "test"
	userId := uint(1)
	exp := 3600

	// JWT 비밀키를 임시로 잘못된 값으로 변경
	originalSecret := config.GetEnv("JWT_SECRET")
	config.SetEnv("JWT_SECRET", "invalid_secret")
	defer config.SetEnv("JWT_SECRET", originalSecret)

	token, err := CreateToken(tokenType, userId, exp)
	assert.NoError(t, err)    // 잘못된 비밀키로 토큰 생성이 성공해야 함
	assert.NotEmpty(t, token) // 토큰이 비어 있지 않아야 함

	// 토큰 검증
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(originalSecret), nil // 원래의 비밀키를 사용해 검증
	})
	assert.Error(t, err) // 잘못된 비밀키로 검증이 실패해야 함
}

// TestCreateTokenSet_InvalidUserId는 유효하지 않은 사용자 ID로 CreateTokenSet 함수가 올바르게 동작하는지 테스트합니다.
func TestCreateTokenSet_InvalidUserId(t *testing.T) {
	invalidUserId := uint(0)

	accessToken, refreshToken, err := CreateTokenSet(invalidUserId)
	assert.NoError(t, err)           // 오류가 없는지 확인
	assert.NotEmpty(t, accessToken)  // 액세스 토큰이 비어있지 않은지 확인
	assert.NotEmpty(t, refreshToken) // 리프레시 토큰이 비어있지 않은지 확인

	// 액세스 토큰 검증
	accessClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "access", accessClaims["token_type"])         // 토큰 타입이 "access"인지 확인
	assert.Equal(t, float64(invalidUserId), accessClaims["user"]) // 사용자 ID가 일치하는지 확인

	// 리프레시 토큰 검증
	refreshClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "refresh", refreshClaims["token_type"])        // 토큰 타입이 "refresh"인지 확인
	assert.Equal(t, float64(invalidUserId), refreshClaims["user"]) // 사용자 ID가 일치하는지 확인
}

// TestCreateTokenSet_ExpiredToken 은 토큰의 만료 시간을 짧게 설정하여 만료된 토큰을 검증합니다.
func TestCreateTokenSet_ExpiredToken(t *testing.T) {
	userId := uint(1)

	// 매우 짧은 만료 시간으로 토큰 생성
	accessToken, err := CreateToken("access", userId, 1)
	refreshToken, err := CreateToken("refresh", userId, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	time.Sleep(2 * time.Second) // 토큰이 만료되도록 2초 대기

	// 만료된 액세스 토큰 검증
	accessClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	assert.Error(t, err) // 만료된 토큰은 오류가 발생해야 함

	// 매우 짧은 만료 시간으로 리프레시 토큰 생성
	refreshToken, err = CreateToken("refresh", userId, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	time.Sleep(2 * time.Second) // 토큰이 만료되도록 2초 대기

	// 만료된 리프레시 토큰 검증
	refreshClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	assert.Error(t, err) // 만료된 토큰은 오류가 발생해야 함
}
