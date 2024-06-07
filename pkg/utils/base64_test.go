package utils

import (
	"reflect"
	"testing"
)

// TestBase64Encode Base64Encode 함수 테스트
func TestBase64Encode(t *testing.T) {
	// 테스트 케이스 목록
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		// 일반 문자열 테스트
		{name: "simple string", input: []byte("hello"), expected: "aGVsbG8="},
		// 빈 문자열 테스트
		{name: "empty string", input: []byte(""), expected: ""},
		// 특수 문자 포함 문자열 테스트
		{name: "special characters", input: []byte("!@#$%^&*()"), expected: "IUAjJCVeJiooKQ=="},
		// 유니코드 문자열 테스트
		{name: "unicode characters", input: []byte("안녕하세요"), expected: "7JWI64WV7ZWY7IS47JqU"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Encode(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestBase64Decode Base64Decode 함수 테스트
func TestBase64Decode(t *testing.T) {
	// 테스트 케이스 목록
	tests := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		// 정상적인 인코딩 문자열 테스트
		{name: "simple string", input: "aGVsbG8=", expected: []byte("hello"), wantErr: false},
		// 빈 문자열 테스트
		{name: "empty string", input: "", expected: []byte(""), wantErr: false},
		// 특수 문자 포함 문자열 테스트
		{name: "special characters", input: "IUAjJCVeJiooKQ==", expected: []byte("!@#$%^&*()"), wantErr: false},
		// 유니코드 문자열 테스트
		{name: "unicode characters", input: "7JWI64WV7ZWY7IS47JqU", expected: []byte("안녕하세요"), wantErr: false},
		// 잘못된 인코딩 문자열 테스트
		{name: "invalid string", input: "invalid!!", expected: nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Base64Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
