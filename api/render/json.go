package render

import (
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

// JSON 상태 코드와 결과 메시지를 설정하여 JSON 형식으로 응답하는 함수입니다.
func JSON(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// message 타입 확인(map, slice, struct)
	if _, ok := message.(map[string]string); ok {
		message = message.(map[string]string)
	} else if _, ok := message.([]map[string]string); ok {
		message = message.([]map[string]string)
	}

	jsonResponse, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}
	_, _ = w.Write(jsonResponse)
}
