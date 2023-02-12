package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
)

func Link(context *fiber.Ctx) error {
	user_id, _ := strconv.Atoi(context.Params("id"))

	var links []models.Link

	database.DB.Where("user_id = ?", user_id).Find(&links)

	for i, link := range links {
		var orders []models.Order

		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return context.JSON(links)
}

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(context *fiber.Ctx) error {
	var request CreateLinkRequest

	if err := context.BodyParser(&request); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(context)

	link := models.Link{
		UserId: id,
		Code:   faker.Username(),
	}

	for _, productId := range request.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)

	return context.JSON(link)
}
