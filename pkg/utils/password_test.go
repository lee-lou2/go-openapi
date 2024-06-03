package utils

import (
	"testing"
)

// TestHashPasswordAndCheckPasswordHash HashPassword 및 CheckPasswordHash 함수 테스트
func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "mysecretpassword"

	// 패스워드 해시 생성 테스트
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// 해시된 패스워드 검증 테스트 (올바른 패스워드)
	if !CheckPasswordHash(password, hash) {
		t.Errorf("CheckPasswordHash() failed for correct password")
	}

	// 해시된 패스워드 검증 테스트 (잘못된 패스워드)
	wrongPassword := "wrongpassword"
	if CheckPasswordHash(wrongPassword, hash) {
		t.Errorf("CheckPasswordHash() succeeded for incorrect password")
	}
}

// TestHashPasswordErrorHandling HashPassword 함수의 에러 처리 테스트
func TestHashPasswordErrorHandling(t *testing.T) {
	// 빈 문자열에 대한 해시 생성 테스트
	_, err := HashPassword("")
	if err != nil {
		t.Errorf("HashPassword() error = %v; expected no error for empty password", err)
	}
}
