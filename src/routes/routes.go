package routes

import (
	"ambassador/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// apiルーティング
	api := app.Group("api")

	// adminルーティング
	admin := api.Group("admin")

	admin.Post("register", controllers.Register)

	admin.Post("login", controllers.Login)
	admin.Post("logout", controllers.Logout)

	admin.Get("user", controllers.GetUser)
}
