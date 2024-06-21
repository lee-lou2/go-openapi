package auth

import (
	"go-openapi/api/render"
	authInternal "go-openapi/internal/auth"
	"log"
	"net/http"
	"strconv"
)

// // CreateClientHandler 클라이언트 키 생성 핸들러
// func CreateClientHandler(c fiber.Ctx) error {
// 	user := fiber.Locals[uint](c, "user")
// 	instance, err := authInternal.CreateClient(user)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(fiber.Map{
// 		"clientId":     instance.ClientId,
// 		"clientSecret": instance.ClientSecret,
// 	})
// }

// CreateClientHandler 클라이언트 키 생성 핸들러
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	instance, err := authInternal.CreateClient(user)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{
		"clientId":     instance.ClientId,
		"clientSecret": instance.ClientSecret,
	})
}

// GetClientsHandler 클라이언트 키 조회 핸들러
func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	clients, err := authInternal.GetClients(user)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		return
	}
	resp := make([]map[string]string, 0)
	for _, instance := range *clients {
		resp = append(resp, map[string]string{
			"id":           strconv.Itoa(int(instance.ID)),
			"clientId":     instance.ClientId,
			"clientSecret": instance.ClientSecret,
			"createdAt":    instance.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	render.JSON(w, http.StatusOK, resp)
}

// DeleteClientHandler 클라이언트 키 삭제 핸들러
func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	// url path parameter 를 조회
	id := r.PathValue("id")
	err := authInternal.DeleteClient(user, id)
	if err != nil {
		log.Println(err.Error())
		render.JSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"message": "Success"})
}
