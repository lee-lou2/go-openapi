package utils

import (
	"testing"
)

// TestAES256Encrypt AES256 Encrypt 함수 테스트
func TestAES256Encrypt(t *testing.T) {
	aes := NewAES256("6368616e676520746869732070617373")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "simple string", input: "hello"},
		{name: "empty string", input: ""},
		{name: "special characters", input: "!@#$%^&*()"},
		{name: "unicode characters", input: "안녕하세요"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted := aes.Encrypt(tt.input)
			if encrypted == "" {
				t.Errorf("expected non-empty string, got empty string")
			}
		})
	}
}

// TestAES256Decrypt AES256 Decrypt 함수 테스트
func TestAES256Decrypt(t *testing.T) {
	aes := NewAES256()

	tests := []struct {
		name  string
		input string
	}{
		{name: "simple string", input: "hello"},
		{name: "empty string", input: ""},
		{name: "special characters", input: "!@#$%^&*()"},
		{name: "unicode characters", input: "안녕하세요"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted := aes.Encrypt(tt.input)
			decrypted := aes.Decrypt(encrypted)
			if decrypted != tt.input {
				t.Errorf("expected %s, got %s", tt.input, decrypted)
			}
		})
	}
}

// TestAES256InvalidDecrypt AES256 잘못된 복호화 테스트
func TestAES256InvalidDecrypt(t *testing.T) {
	aes := NewAES256()

	invalidInput := "invalid_encrypted_string"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()

	aes.Decrypt(invalidInput)
}
