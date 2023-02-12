package main

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/go-redis/redis/v8"
)

// 初期データの作成
func main() {
	database.Connect()
	CreateUsers()
	CreateProducts()
	CreateOrders()
	CreateRankings()
}

func CreateUsers() {
	for i := 0; i < 30; i++ {
		ambassador := models.User{
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}

		ambassador.SetUserPassword("1234")
		database.DB.Create(&ambassador)
	}
}

func CreateProducts() {
	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Username(),
			Description: faker.Username(),
			Image:       faker.URL(),
			Price:       float64(rand.Intn(90) + 10),
		}

		database.DB.Create(&product)
	}
}

func CreateOrders() {
	for i := 0; i < 30; i++ {
		var orderItems []models.OrderItem
		for j := 0; j < rand.Intn(5); j++ {
			price := float64(rand.Intn(90) + 10)
			quantity := uint(rand.Intn(5))

			orderItems = append(orderItems, models.OrderItem{
				ProductTitle:      faker.Word(),
				Price:             price,
				Quantity:          quantity,
				AdminRevenue:      0.9 * price * float64(quantity),
				AmbassadorRevenue: 0.1 * price * float64(quantity),
			})
		}

		database.DB.Create(&models.Order{
			UserId:          uint(rand.Intn(30) + 1),
			Code:            faker.Username(),
			AmbassadorEmail: faker.Email(),
			FirstName:       faker.FirstName(),
			LastName:        faker.LastName(),
			Email:           faker.Email(),
			Complete:        true,
			OrderItems:      orderItems,
		})
	}
}

func CreateRankings() {
	database.SetupRedis()

	ctx := context.Background()

	var users []models.User

	database.DB.Find(&users, models.User{
		IsAmbassador: true,
	})

	for _, user := range users {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(database.DB)

		// スコアをつける
		database.Cache.ZAdd(ctx, "rankings", &redis.Z{
			Score:  *ambassador.Revenue,
			Member: user.Name(),
		})
	}
}
