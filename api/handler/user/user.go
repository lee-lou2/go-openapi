package user

import (
	"go-openapi/api/response"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
	"net/http"
)

// CreateUserHandler 사용자 생성 핸들러
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if !validation.ValidateEmail(email) || !validation.ValidatePassword(password) {
		response.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	if err := userInternal.CreateUser(email, password); err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	response.JSON(w, http.StatusCreated, map[string]string{"email": email})
}
