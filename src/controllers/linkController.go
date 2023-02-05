package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Link(context *fiber.Ctx) error {
	user_id, _ := strconv.Atoi(context.Params("id"))

	var links []models.Link

	database.DB.Where("user_id = ?", user_id).Find(&links)

	return context.JSON(links)
}
