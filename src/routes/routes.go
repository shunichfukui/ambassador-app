package routes

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// apiルーティング
	api := app.Group("api")

	// adminルーティング
	admin := api.Group("admin")

	admin.Post("register", controllers.Register)

	admin.Post("login", controllers.Login)

	adminAuthenticated := admin.Use(middlewares.IsAuthenticated)
	adminAuthenticated.Post("logout", controllers.Logout)
	adminAuthenticated.Get("user", controllers.GetUser)
}
