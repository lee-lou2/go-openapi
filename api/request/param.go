package request

import "net/http"

// GetQueryParams 요청의 쿼리 파라미터 조회
func GetQueryParams(r *http.Request) map[string]string {
	query := r.URL.Query()
	result := make(map[string]string)
	for key, value := range query {
		result[key] = value[0]
	}
	return result
}
