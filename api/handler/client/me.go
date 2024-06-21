package client

import (
	"go-openapi/api/render"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GetMeHandler 내 정보 조회 핸들러
func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	clientId := r.Context().Value("client").(uint)
	render.JSON(w, fiber.StatusOK, map[string]string{
		"id": strconv.Itoa(int(clientId)),
	})
}
