package database

import (
	"ambassador/src/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	user := os.Getenv("MySQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbname)

	DB, err = gorm.Open(mysql.Open(connection), &gorm.Config{})

	if err != nil {
		panic("データベースの接続に失敗しました!")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{}, models.Product{}, models.Link{}, models.Order{}, models.OrderItem{})
}
