package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(context *fiber.Ctx) error {
	var data map[string]string

	if err := context.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		context.Status(400)

		return context.JSON(fiber.Map{
			"message": "最初に入力したパスワードと確認パスワードの値が正しくありません。",
		})
	}

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(context.Path(), "/api/ambassador"),
	}

	user.SetUserPassword(data["password"])

	database.DB.Create(&user)

	return context.JSON(user)
}

func Login(context *fiber.Ctx) error {
	var data map[string]string

	if err := context.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": "ユーザーが見つかりませんでした",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": "パスワードが違います。",
		})
	}

	isAmbassador := strings.Contains(context.Path(), "/api/ambassador")

	var scope string

	if isAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	if !isAmbassador && user.IsAmbassador {
		context.Status(fiber.StatusUnauthorized)
		return context.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	token, err := middlewares.GenerateJWT(user.Id, scope)

	if err != nil {
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	context.Cookie(&cookie)

	return context.JSON(fiber.Map{
		"message": "success",
	})
}

// remove cookie
func Logout(context *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	context.Cookie(&cookie)

	return context.JSON(fiber.Map{
		"message": "success",
	})
}

func GetUser(context *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(context)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	if strings.Contains(context.Path(), "/api/ambassador") {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(database.DB)
		return context.JSON(ambassador)
	}

	return context.JSON(user)
}

func UpdateUserInfo(context *fiber.Ctx) error {
	var data map[string]string

	if err := context.BodyParser(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(context)

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(context.Path(), "/api/ambassador"),
	}

	user.Id = id

	database.DB.Model(&user).Updates(&user)

	return context.JSON(user)
}

func UpdateUserPassword(context *fiber.Ctx) error {
	var data map[string]string

	if err := context.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		context.Status(400)

		return context.JSON(fiber.Map{
			"message": "最初に入力したパスワードと確認パスワードの値が正しくありません。",
		})
	}

	id, _ := middlewares.GetUserId(context)

	user := models.User{}
	user.Id = id

	user.SetUserPassword(data["password"])

	database.DB.Model(&user).Updates(&user)

	return context.JSON(user)
}
