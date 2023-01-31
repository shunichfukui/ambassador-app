package routes

import (
	"ambassador/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// apiルーティング
	api := app.Group("api")
	api.Post("/admin/register", controllers.Register)

	// adminルーティング
	admin := api.Group("admin")
	admin.Post("/admin/register", controllers.Register)

	app.Get("/", func(context *fiber.Ctx) error {
		return context.SendString("Hello, World 👋!")
	})
}
