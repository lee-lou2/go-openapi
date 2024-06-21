package user

import (
	"go-openapi/api/parser"
	"go-openapi/api/render"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
	"net/http"
)

// SendPasswordResetCodeHandler 비밀번호 재설정 코드 전송 핸들러
func SendPasswordResetCodeHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := userInternal.SendPasswordVerifyCode(req.Email); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"message": "Code sent"})
}

// ResetPasswordHandler 비밀번호 재설정 핸들러
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code     string `uri:"code"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := parser.JSON(r, &req); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	if !validation.ValidateEmail(req.Email) || !validation.ValidatePassword(req.Password) || !validation.ValidateCode(req.Code) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	if err := userInternal.VerifyCodeAndChangePassword(req.Email, req.Code, req.Password); err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"message": "Password updated"})
}
