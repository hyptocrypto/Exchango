package main

import (
	"github.com/hyptocrypto/go_exchange_api/server/config"
	"github.com/hyptocrypto/go_exchange_api/server/models"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup(db *gorm.DB) {
	db.AutoMigrate(&models.Trading_Pairs{}, &models.Orders{})
	seed(db)
}

func seed(db *gorm.DB) {
	trading_pairs := []models.Trading_Pairs{
		{Ticker: "USDT_BTC", Price: 40000, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Percent_Change: 0},
		{Ticker: "USDT_LTC", Price: 150, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Percent_Change: 0},
		{Ticker: "USDT_ETH", Price: 1200, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Percent_Change: 0},
		{Ticker: "USDT_XMR", Price: 170, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Percent_Change: 0},
		{Ticker: "USDT_DASH", Price: 100, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Percent_Change: 0},
	}
	for _, pair := range trading_pairs {
		db.Create(&pair)
	}
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

	setup(db)
}
