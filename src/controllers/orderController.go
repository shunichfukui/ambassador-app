package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

func GetOrders(context *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}

	return context.JSON(orders)
}

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrders(context *fiber.Ctx) error {
	var request CreateOrderRequest

	if err := context.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{
		Code: request.Code,
	}

	database.DB.Preload("User").First(&link)

	if link.Id == 0 {
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": "リンクが有効ではありません。",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	// transaction
	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		orderItem := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		// transaction
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			context.Status(fiber.StatusBadRequest)
			return context.JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(requestProduct["quantity"])),
		})
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("http://localhost:500/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:500/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}

	source, err := session.New(&params)

	if err != nil {
		tx.Rollback()
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	order.TransactionId = source.ID

	// transaction
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		context.Status(fiber.StatusBadRequest)
		return context.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tx.Commit()

	return context.JSON(source)
}

func CompleteOrder(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	order := models.Order{}

	database.DB.Preload("OrderItems").First(&order, models.Order{
		TransactionId: data["source"],
	})

	if order.Id == 0 {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"message": "該当の注文が見つかりませんでした",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	go func(order models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			ambassadorRevenue += item.AmbassadorRevenue
			adminRevenue += item.AdminRevenue
		}

		user := models.User{}
		user.Id = order.UserId

		database.DB.First(&user)

		// ランキング更新
		database.Cache.ZIncrBy(context.Background(), "rankings", ambassadorRevenue, user.Name())

		ambassadorMessage := []byte(fmt.Sprintf("you earned $%f from the link #%s", ambassadorRevenue, order.Code))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{order.AmbassadorEmail}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("order #%d with a total of $%f", order.Id, adminRevenue))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
	}(order)

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
