package user

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
	"log"
	"net/http"
)

// CreateUserHandler 사용자 생성 핸들러
func CreateUserHandler(c fiber.Ctx) error {
	body := new(struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	})
	if err := c.Bind().Body(body); err != nil {
		return err
	}
	if !validation.ValidateEmail(body.Email) || !validation.ValidatePassword(body.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	err := userInternal.CreateUser(body.Email, body.Password)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"email": body.Email,
	})
}

// WriteJSONResponse 상태 코드와 결과 메시지를 설정하여 JSON 형식으로 응답하는 함수입니다.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, message map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResponse, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}
	w.Write(jsonResponse)
}

// CreateUserHandler2 사용자 생성 핸들러
func CreateUserHandler2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"message": "Method not allowed"})
		return
	}
	if err := r.ParseForm(); err != nil {
		WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if !validation.ValidateEmail(email) || !validation.ValidatePassword(password) {
		WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	err := userInternal.CreateUser(email, password)
	if err != nil {
		WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
		return
	}
	WriteJSONResponse(w, http.StatusCreated, map[string]string{"email": email})
}
