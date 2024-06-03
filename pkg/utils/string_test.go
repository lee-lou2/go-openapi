package utils

import (
	"testing"
)

// TestGenerateRandomString GenerateRandomString 함수 테스트
func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		length int
	}{
		{length: 10},
		{length: 20},
		{length: 0},
		{length: 50},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := GenerateRandomString(tt.length)
			if err != nil {
				t.Fatalf("GenerateRandomString(%d) error = %v", tt.length, err)
			}
			if len(result) != tt.length {
				t.Errorf("GenerateRandomString(%d) length = %d; expected %d", tt.length, len(result), tt.length)
			}
			for _, char := range result {
				if !contains(letters, char) {
					t.Errorf("GenerateRandomString(%d) contains invalid character %c", tt.length, char)
				}
			}
		})
	}
}

// contains 함수는 문자열 내에 특정 문자가 있는지 확인합니다.
func contains(s string, c rune) bool {
	for _, char := range s {
		if char == c {
			return true
		}
	}
	return false
}
