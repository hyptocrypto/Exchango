package main

import (
	"fmt"

	"github.com/hyptocrypto/go_exchange_api/server/config"
	"github.com/hyptocrypto/go_exchange_api/server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func select_first_order(db *gorm.DB) []models.Orders {
	var orders []models.Orders
	db.Order("created_at asc").Find(&orders)
	orders.
	for _, order := range orders {
		fmt.Println(order.CreatedAt)
		fmt.Println("-------------------------------")
	}
	return orders
}

func main() {
	db, err := gorm.Open(sqlite.Open(config.DB_Path), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, sqlite_err := db.DB()
	if sqlite_err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	select_first_order(db)
}
