package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"
	"context"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/go-redis/redis/v8"
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

func Stats(context *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(context)

	var links []models.Link

	database.DB.Find(&links, models.Link{
		UserId: id,
	})

	var result []interface{}

	var orders []models.Order

	for _, link := range links {
		database.DB.Preload("OrderItems").Find(&orders, &models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0

		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}

	return context.JSON(result)
}

func Rankings(c *fiber.Ctx) error {
	// 並び替え
	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return err
	}

	result := make(map[string]float64)

	for _, ranking := range rankings {
		result[ranking.Member.(string)] = ranking.Score
	}

	return c.JSON(result)
}

func GetLink(context *fiber.Ctx) error {
	code := context.Params("code")
	link := models.Link{
		Code: code,
	}

	database.DB.Preload("User").Preload("Products").First(&link)

	return context.JSON(link)
}
