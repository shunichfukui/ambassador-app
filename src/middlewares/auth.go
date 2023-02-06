package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const SecretKey = "secret"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(context *fiber.Ctx) error {
	cookie := context.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		context.Status(fiber.StatusUnauthorized)
		return context.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)
	IsAmbassador := strings.Contains(context.Path(), "/api/ambassador")

	if (payload.Scope == "admin" && IsAmbassador) || (payload.Scope == "ambassador" && !IsAmbassador) {
		context.Status(fiber.StatusUnauthorized)
		return context.JSON(fiber.Map{
			"message": "認証に失敗しました。",
		})
	}

	return context.Next()
}

func GenerateJWT(id uint, scope string) (string, error) {
	payload := ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecretKey))
}

func GetUserId(context *fiber.Ctx) (uint, error) {
	cookie := context.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*ClaimsWithScope)

	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil
}
