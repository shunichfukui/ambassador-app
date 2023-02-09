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
	adminAuthenticated.Put("user/info", controllers.UpdateUserInfo)
	adminAuthenticated.Put("user/password", controllers.UpdateUserPassword)
	adminAuthenticated.Get("ambassadors", controllers.GetAmbassadors)

	// products
	adminAuthenticated.Get("products", controllers.GetProducts)
	adminAuthenticated.Post("products", controllers.CreateProducts)
	adminAuthenticated.Get("products/:id", controllers.GetProduct)
	adminAuthenticated.Put("products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("products/:id", controllers.DeleteProduct)

	// link
	adminAuthenticated.Delete("users/:id/links", controllers.Link)

	// orders
	adminAuthenticated.Get("orders", controllers.GetOrders)

	// ambassadorルーティング
	ambassador := api.Group("ambassador")
	ambassador.Post("register", controllers.Register)
	ambassador.Post("login", controllers.Login)
	ambassador.Post("products/frontend", controllers.ProductsFrontend)

	ambassadorAuthenticated := ambassador.Use(middlewares.IsAuthenticated)
	ambassadorAuthenticated.Post("logout", controllers.Logout)
	ambassadorAuthenticated.Get("user", controllers.GetUser)
	ambassadorAuthenticated.Put("user/info", controllers.UpdateUserInfo)
	ambassadorAuthenticated.Put("user/password", controllers.UpdateUserPassword)
}
