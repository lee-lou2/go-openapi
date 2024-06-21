package middleware

import (
	"context"
	"fmt"
	"go-openapi/api/render"
	"go-openapi/config"
	clientModel "go-openapi/model/client"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
)

var ctx = context.Background()

// LimitPerSecondMiddleware 초당 요청 제한 미들웨어
func LimitPerSecondMiddleware(scope string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientId := r.Context().Value("client").(uint)
		limit := clientModel.ScopeRateLimits[scope]

		// 캐시를 이용해서 요청 수 체크
		cache := config.GetCache()
		key := fmt.Sprintf("%s:%d", scope, clientId)
		count, err := cache.Incr(ctx, key).Result()
		if err != nil {
			render.JSON(w, fiber.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		}
		if count == 1 {
			cache.Expire(ctx, key, time.Second)
		}
		if count > int64(limit) {
			db := config.GetDB()
			var result struct {
				Level int
			}
			db.Table("clients").Select("users.level").Joins("left join users on users.id = clients.user_id").Where("clients.id = ?", clientId).Scan(&result)
			switch result.Level {
			case 0:
				render.JSON(w, fiber.StatusTooManyRequests, map[string]string{"error": "Too Many Requests"})
				return
			case 1:
				premiumLimit := limit * 2
				if count > int64(premiumLimit) {
					render.JSON(w, fiber.StatusTooManyRequests, map[string]string{"error": "Too Many Requests"})
					return
				}
			case 2:
				vipLimit := limit * 3
				if count > int64(vipLimit) {
					render.JSON(w, fiber.StatusTooManyRequests, map[string]string{"error": "Too Many Requests"})
					return
				}
			default:
				render.JSON(w, fiber.StatusTooManyRequests, map[string]string{"error": "Too Many Requests"})
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
