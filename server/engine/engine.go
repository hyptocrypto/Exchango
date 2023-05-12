package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/hyptocrypto/go_exchange_api/server/models"
	"github.com/hyptocrypto/go_exchange_api/server/utils"
	"gorm.io/gorm"
)

func update_data_live(db *gorm.DB) {
	for {
		update_data(db)
		time.Sleep(60 * time.Second)
	}
}
func update_worker(db *gorm.DB, data interface{}, pair string, wg *sync.WaitGroup) {
	defer wg.Done()
	db.Model(&models.Trading_Pairs{}).Where("Ticker = ?", pair).Updates(models.Trading_Pairs{
		Price:          utils.Interface_to_float(data.(map[string]interface{})["last"]),
		Daily_Volume:   utils.Interface_to_float(data.(map[string]interface{})["baseVolume"]),
		Daily_High:     utils.Interface_to_float(data.(map[string]interface{})["high24hr"]),
		Daily_Low:      utils.Interface_to_float(data.(map[string]interface{})["low24hr"]),
		Percent_Change: utils.Interface_to_float(data.(map[string]interface{})["percentChange"]),
	})

}

func update_data(db *gorm.DB) {
	resp, http_error := http.Get("https://poloniex.com/public?command=returnTicker")
	if http_error != nil {
		panic(http_error)
	}
	body, read_error := ioutil.ReadAll(resp.Body)
	if read_error != nil {
		panic(read_error)
	}
	var data map[string]interface{}

	data_json_error := json.Unmarshal([]byte(string(body)), &data)
	if data_json_error != nil {
		panic(data_json_error)
	}

	pairs := map[string]interface{}{
		"USDT_BTC":  data["USDT_BTC"],
		"USDT_ETH":  data["USDT_ETH"],
		"USDT_XMR":  data["USDT_XMR"],
		"USDT_LTC":  data["USDT_LTC"],
		"USDT_DASH": data["USDT_DASH"]}
	var wg sync.WaitGroup

	for key, value := range pairs {
		wg.Add(1)
		go update_worker(db, value, key, &wg)
	}
	wg.Wait()
}

func open_sell_orders(db *gorm.DB) []models.Orders {
	var all_data []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Price in (select price from Orders group by price having count(*)>1) AND Order_Type=? AND Settled=?", "Sell", false)
	// db.Raw("SELECT * FROM models.Orders WHERE Price in (select price from models.Orders group by price having count(*)>1) AND Order_type=?", "Sell").Scan(&all_data)
	return all_data
}

func open_buy_orders(db *gorm.DB) []models.Orders {
	var all_data []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Price in (select price from Orders group by price having count(*)>1) AND Order_Type=? And Settled=?", "Buy", false)
	return all_data
}

func Settle_orders_live(db *gorm.DB) {
	for {
		settle_orders(db)
		time.Sleep(1 * time.Second)

	}
}
func settle_orders(db *gorm.DB) {
	buy_orders := open_buy_orders(db)
	sell_orders := open_sell_orders(db)
	for buy_index, buy := range buy_orders {
		for sell_index, sell := range sell_orders {
			fmt.Println(buy_index, sell_index)
			fmt.Println(buy.Trading_Pair.Ticker, sell.Trading_Pair.Ticker)
			if buy.Trading_Pair.Ticker != sell.Trading_Pair.Ticker {
			} else {
				fmt.Println("buy", buy_orders[buy_index].Current_Amount)
				fmt.Println("sell", sell_orders[sell_index].Current_Amount)
				if buy_orders[buy_index].Current_Amount > sell_orders[sell_index].Current_Amount && sell_orders[sell_index].Current_Amount > 0 {
					fmt.Println("buy - sell")
					fmt.Println(buy_orders[buy_index].Current_Amount, sell_orders[sell_index].Current_Amount)
					buy_orders[buy_index].Current_Amount = buy.Current_Amount - sell.Current_Amount
					sell_orders[sell_index].Current_Amount = 0
					fmt.Println("buy new value ", buy_orders[buy_index].Current_Amount)
					fmt.Println("sell new value", sell_orders[sell_index].Current_Amount)

					var buy_order models.Orders
					var sell_order models.Orders
					db.Find(&buy_order, "ID=?", buy.ID)
					db.Model(&buy_order).Updates(map[string]interface{}{"Current_Amount": buy_orders[buy_index].Current_Amount, "Partial_Settled": true})
					db.Find(&sell_order, "ID=?", sell.ID)
					db.Model(&sell_order).Updates(map[string]interface{}{"Current_Amount": 0, "Settled": true})

				}
				if sell_orders[sell_index].Current_Amount > buy_orders[buy_index].Current_Amount && buy_orders[buy_index].Current_Amount > 0 {
					fmt.Println("sell - buy")
					sell_orders[sell_index].Current_Amount = sell_orders[sell_index].Current_Amount - buy_orders[sell_index].Current_Amount
					buy_orders[buy_index].Current_Amount = 0
					fmt.Println(sell_orders[sell_index].Current_Amount, buy_orders[buy_index].Current_Amount)
					fmt.Println("sell new value ", sell_orders[sell_index].Current_Amount)

					var buy_order models.Orders
					var sell_order models.Orders
					db.Find(&sell_order, "ID=?", sell.ID)
					db.Model(&sell_order).Updates(map[string]interface{}{"Current_Amount": sell_orders[sell_index].Current_Amount, "Partial_Settled": true})
					db.Find(&buy_order, "ID=?", buy.ID)
					db.Model(&buy_order).Updates(map[string]interface{}{"Current_Amount": 0, "Settled": true})

				}
				if buy_orders[buy_index].Current_Amount == sell_orders[sell_index].Current_Amount {
					fmt.Println("buy = sell")
					fmt.Println(buy_orders[buy_index].Current_Amount, sell_orders[sell_index].Current_Amount)
					buy_orders[buy_index].Current_Amount = 0
					sell_orders[sell_index].Current_Amount = 0
					fmt.Println("new value ", buy_orders[buy_index].Current_Amount)

					var buy_order models.Orders
					var sell_order models.Orders
					db.Find(&buy_order, "ID=?", buy.ID)
					db.Model(&buy_order).Updates(map[string]interface{}{"Current_Amount": 0, "Settled": true})
					db.Find(&sell_order, "ID=?", sell.ID)
					db.Model(&sell_order).Updates(map[string]interface{}{"Current_Amount": 0, "Settled": true})

				}

			}
		}

	}
}
