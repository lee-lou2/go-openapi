package utils

import (
	"testing"
)

// TestSHA256 다양한 테스트 케이스로 SHA256 함수를 검증
func TestSHA256(t *testing.T) {
	// 테스트 케이스 목록
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			input:    "world",
			expected: "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7",
		},
		{
			input:    "golang",
			expected: "d754ed9f64ac293b10268157f283ee23256fb32a4f8dedb25c8446ca5bcb0bb3",
		},
		{
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			input:    "1234567890",
			expected: "c775e7b757ede630cd0aa1113bd102661ab38829ca52a6422ab782862f268646",
		},
	}

	for _, tc := range testCases {
		// SHA256 함수 호출
		result := SHA256(tc.input)
		// 결과 검증
		if result != tc.expected {
			t.Errorf("SHA256(%s) = %s; expected %s", tc.input, result, tc.expected)
		}
	}
}
