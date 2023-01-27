package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	DB, err = gorm.Open(mysql.Open("root:root@tcp(ambassador_mysql)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("データベースの接続に失敗しました!")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{})
}
