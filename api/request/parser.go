package request

import (
	"github.com/goccy/go-json"
	"net/http"
)

// ParseJSONBody 함수는 요청의 JSON 바디를 파싱하여 v에 저장
func ParseJSONBody(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}
