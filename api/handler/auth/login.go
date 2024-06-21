package auth

import (
	"fmt"
	"go-openapi/api/render"
	"go-openapi/api/validation"
	authInternal "go-openapi/internal/auth"
	"log"
	"net/http"
)

// LoginHandler 로그인 핸들러
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if !validation.ValidateEmail(email) || !validation.ValidatePassword(password) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	accessToken, refreshToken, err := authInternal.GetTokenFromLogin(email, password)
	if err != nil {
		// 모든 오류 내용 통일
		log.Println(err.Error())
		render.JSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{
		"tokenType":             "Bearer",
		"accessToken":           accessToken,
		"refreshToken":          refreshToken,
		"accessTokenExpiresIn":  "3600",
		"refreshTokenExpiresIn": "86400",
	})
}

// LogoutHandler 로그아웃 핸들러
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	fmt.Println(user)
	render.JSON(w, http.StatusOK, map[string]string{"message": "logout"})
}
