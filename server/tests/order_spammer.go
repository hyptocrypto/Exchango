package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hyptocrypto/go_exchange_api/server/config"
	"github.com/hyptocrypto/go_exchange_api/server/models"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// func (db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		var data map[string]interface{}
// 		_ = json.NewDecoder(r.Body).Decode(&data)
// 		amount_val := interface_to_float(data["Amount"].(string))
// 		price_val := interface_to_float(data["Price"].(string))
// 		settled_val := string_to_bool(data["Settled"].(string))
// 		var trading_pair Trading_Pair
// 		db.First(&trading_pair, "Ticker=?", data["Trading_Pair"])

// 		order := Orders{Trading_PairID: trading_pair.ID,
// 			Trading_Pair:    trading_pair,
// 			Order_Type:      data["Order_type"].(string),
// 			Opening_Amount:  amount_val,
// 			Current_Amount:  amount_val,
// 			Price:           price_val,
// 			Settled:         settled_val,
// 			Partial_Settled: settled_val}
// 		db.Create(&order)

// 	}
// }

var prices = [...]int{40000, 45000, 30000, 20000, 35000, 25000}
var amounts = [...]int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50}
var trading_pairs = map[int]string{1: "USDT_BTC", 3: "USDT_ETH", 2: "USDT_LTC", 4: "USDT_XMR", 5: "USDT_DASH"}
var types = [...]string{"Buy", "Sell"}

// initialize global pseudo random generator
// message := fmt.Sprint("Gonna work from home...", reasons[rand.Intn(len(reasons))])

func create_orders(db *gorm.DB) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 50; i++ {
		a := randomInt(1, 6)
		amount := float64(amounts[rand.Intn(len(amounts))])
		var pair models.Trading_Pairs
		db.First(&pair, "Ticker=?", trading_pairs[a])
		order := models.Orders{Trading_PairID: uint(a),
			Trading_Pair:    pair,
			Order_Type:      types[rand.Intn(len(types))],
			Opening_Amount:  amount,
			Current_Amount:  amount,
			Price:           float64(prices[rand.Intn(len(prices))]),
			Settled:         false,
			Partial_Settled: false}
		db.Create(&order)
		fmt.Println(order)
		time.Sleep(2 * time.Second)
	}

}
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
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

	fmt.Println("Starting Order Creation")
	create_orders(db)
}
