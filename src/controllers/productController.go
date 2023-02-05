package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"

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