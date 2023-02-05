package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
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
		IsAmbassador: false,
	}

	password := user.SetUserPassword(data["password"])

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

	var payload := jwt.StandardClaims{
		Subject: strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix()
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))

	if err != nil {
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24),
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
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	context.Cookie(&cookie)

	return context.JSON(fiber.Map{
		"message": "success",
	})
}

func GetUser(context *fiber.Ctx) error {
	cookie := context.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		context.Status(fiber.StatusUnauthorized)
		return context.JSON(fiber.Map{
			"message": "認証ができませんでした。",
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", payload.Subject).First(&user)

	return context.JSON(user)
}