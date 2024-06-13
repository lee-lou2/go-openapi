package user

import (
	"go-openapi/api/request"
	"go-openapi/api/response"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
	"net/http"
)

// SendVerifyCodeHandler 인증 코드 전송 핸들러
func SendVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := request.ParseJSONBody(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	if !validation.ValidateEmail(req.Email) {
		response.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	if err := userInternal.ValidateUserAndSendVerifyCode(req.Email); err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Code sent"})
}

// VerifyCodeHandler 인증 코드 검증 핸들러
func VerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := request.GetPathParam(r, "code")
	var req struct {
		Email string `json:"email"`
	}
	if err := request.ParseJSONBody(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	if !validation.ValidateEmail(req.Email) || !validation.ValidateCode(code) {
		response.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	if err := userInternal.VerifyCodeAndUpdateUser(req.Email, code); err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "User verified"})
}
