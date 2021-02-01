package main

import (
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type trading_pair struct {
	gorm.Model
	Ticker         string
	Price          float64
	Daily_Volume   float64
	Daily_High     float64
	Daily_Low      float64
	Precent_Change float64
}

type orders struct {
	gorm.Model
	Trading_PairID  uint
	Trading_Pair    trading_pair
	Order_Type      string
	Opening_Amount  float64
	Current_Amount  float64
	Settled         bool
	Partial_Settled bool
	Price           float64
}

func setup(db *gorm.DB) {
	db.AutoMigrate(&trading_pair{}, &orders{})
	seed(db)
}

func seed(db *gorm.DB) {
	trading_pairs := []trading_pair{
		{Ticker: "USDT_BTC", Price: 40000, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Precent_Change: 0},
		{Ticker: "USDT_LTC", Price: 150, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Precent_Change: 0},
		{Ticker: "USDT_ETH", Price: 1200, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Precent_Change: 0},
		{Ticker: "USDT_XMR", Price: 170, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Precent_Change: 0},
		{Ticker: "USDT_DASH", Price: 100, Daily_Volume: 0, Daily_High: 0, Daily_Low: 0, Precent_Change: 0},
	}
	for _, pair := range trading_pairs {
		db.Create(&pair)
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("mock_exchange.db"), &gorm.Config{})
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
