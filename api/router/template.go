package router

import "github.com/gofiber/fiber/v3"

// TemplateRouter 템플릿 라우터
func TemplateRouter(app fiber.Router) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Get("/accounts/login", func(c fiber.Ctx) error {
		return c.Render("accounts/login", fiber.Map{})
	})
	app.Get("/accounts/signup", func(c fiber.Ctx) error {
		return c.Render("accounts/signup", fiber.Map{})
	})
	app.Get("/accounts/password", func(c fiber.Ctx) error {
		return c.Render("accounts/password", fiber.Map{})
	})
	app.Get("/verify/:codeType", func(c fiber.Ctx) error {
		codeType := fiber.Params[int](c, "codeType")
		if codeType == 1 {
			return c.Render("user/verify", fiber.Map{})
		} else if codeType == 2 {
			return c.Render("user/password", fiber.Map{})
		}
		return c.Redirect().To("/error?m=Invalid code type")
	})
	app.Get("/error", func(c fiber.Ctx) error {
		message := c.Query("m")
		return c.Render("error", fiber.Map{
			"message": message,
		})
	})
}
