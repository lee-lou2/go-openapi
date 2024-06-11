package middleware

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-openapi/config"
	clientModel "go-openapi/model/client"
	"time"
)

var ctx = context.Background()

// LimitPerSecondMiddleware 초당 요청 제한 미들웨어
func LimitPerSecondMiddleware(scope string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		clientId := fiber.Locals[uint](c, "client")
		limit := clientModel.ScopeRateLimits[scope]

		// 캐시를 이용해서 요청 수 체크
		cache := config.GetCache()
		key := fmt.Sprintf("%s:%d", scope, clientId)
		count, err := cache.Incr(ctx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
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
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too Many Requests",
				})
			case 1:
				premiumLimit := limit * 2
				if count > int64(premiumLimit) {
					return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
						"error": "Too Many Requests",
					})
				}
			case 2:
				vipLimit := limit * 3
				if count > int64(vipLimit) {
					return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
						"error": "Too Many Requests",
					})
				}
			default:
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too Many Requests",
				})
			}
		}
		return c.Next()
	}
}
