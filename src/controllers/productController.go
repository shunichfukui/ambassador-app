package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetProducts(context *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)

	return context.JSON(products)
}

func CreateProducts(context *fiber.Ctx) error {
	var product models.Product
	if err := context.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	return context.JSON(product)
}

func GetProduct(context *fiber.Ctx) error {
	var product models.Product
	id, _ := strconv.Atoi(context.Params("id"))

	product.Id = uint(id)

	database.DB.Find(&product)

	return context.JSON(product)
}

func UpdateProduct(context *fiber.Ctx) error {
	id, _ := strconv.Atoi(context.Params("id"))

	product := models.Product{}
	product.Id = uint(id)

	if err := context.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	return context.JSON(product)
}

func DeleteProduct(context *fiber.Ctx) error {
	id, _ := strconv.Atoi(context.Params("id"))

	product := models.Product{}
	product.Id = uint(id)

	database.DB.Delete(&product)

	return nil
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_frontend").Result()

	if err != nil {
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		// キャッシュのセット
		if errKey := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); err != nil {
			panic(errKey)
		}
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	return c.JSON(products)
}
