package user

import (
	"github.com/gofiber/fiber/v3"
	"go-openapi/api/request"
	"go-openapi/api/response"
	"go-openapi/api/validation"
	userInternal "go-openapi/internal/user"
	"net/http"
)

// VerifyCodeHandler 인증 코드 검증 핸들러
func VerifyCodeHandler(c fiber.Ctx) error {
	code := fiber.Params[string](c, "code")
	body := new(struct {
		Email string `json:"email"`
	})
	if err := c.Bind().JSON(body); err != nil {
		return err
	}
	if !validation.ValidateEmail(body.Email) || !validation.ValidateCode(code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})

	}
	err := userInternal.VerifyCodeAndUpdateUser(body.Email, code)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "User verified"})
}

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
