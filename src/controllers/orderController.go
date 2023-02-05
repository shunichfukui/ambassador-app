package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetOrders(context *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Find(&orders)

	return context.JSON(orders)
}
