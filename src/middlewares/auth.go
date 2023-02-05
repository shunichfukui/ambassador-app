package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func IsAuthenticated(context *fiber.Ctx) error {
	cookie := context.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		context.Status(fiber.StatusUnauthorized)
		return context.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	return context.Next()
}

func GetUserId(context *fiber.Ctx) (uint, error) {
	cookie := context.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*jwt.StandardClaims)

	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil
}
