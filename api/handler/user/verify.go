package user

import (
	"go-openapi/api/parser"
	"go-openapi/api/render"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
	"net/http"
)

// SendVerifyCodeHandler 인증 코드 전송 핸들러
func SendVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := parser.JSON(r, &req); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	if !validation.ValidateEmail(req.Email) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	if err := userInternal.ValidateUserAndSendVerifyCode(req.Email); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"message": "Code sent"})
}

// VerifyCodeHandler 인증 코드 검증 핸들러
func VerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var req struct {
		Email string `json:"email"`
	}
	if err := parser.JSON(r, &req); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	if !validation.ValidateEmail(req.Email) || !validation.ValidateCode(code) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	if err := userInternal.VerifyCodeAndUpdateUser(req.Email, code); err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"message": "User verified"})
}
