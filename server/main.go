package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/hyptocrypto/go_exchange_api/server/config"
	"github.com/hyptocrypto/go_exchange_api/server/models"
	"github.com/hyptocrypto/go_exchange_api/server/websocket"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
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
		Price:          interface_to_float(data.(map[string]interface{})["last"]),
		Daily_Volume:   interface_to_float(data.(map[string]interface{})["baseVolume"]),
		Daily_High:     interface_to_float(data.(map[string]interface{})["high24hr"]),
		Daily_Low:      interface_to_float(data.(map[string]interface{})["low24hr"]),
		Percent_Change: interface_to_float(data.(map[string]interface{})["percentChange"]),
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

func interface_to_float(data interface{}) float64 {
	f, err := strconv.ParseFloat(data.(string), 64)
	if err != nil {
		panic(err)
	}
	return f
}
func string_to_bool(data interface{}) bool {
	f, err := strconv.ParseBool(data.(string))
	if err != nil {
		panic(err)
	}
	return f
}

func get_all_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []models.Trading_Pairs
		db.Find(&all_data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}

func get_btc_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_BTC").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func get_eth_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_ETH").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func get_xmr_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_XMR").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func get_ltc_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_LTC").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func get_dash_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_DASH").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func new_order(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&data)
		amount_val := interface_to_float(data["Amount"].(string))
		price_val := interface_to_float(data["Price"].(string))
		settled_val := string_to_bool(data["Settled"].(string))
		var trading_pair models.Trading_Pairs
		db.First(&trading_pair, "Ticker=?", data["Trading_Pair"])
		order := models.Orders{Trading_PairID: trading_pair.ID,
			Trading_Pair:    trading_pair,
			Order_Type:      data["Order_type"].(string),
			Opening_Amount:  amount_val,
			Current_Amount:  amount_val,
			Price:           price_val,
			Settled:         settled_val,
			Partial_Settled: settled_val}
		db.Create(&order)

	}
}
func get_open_orders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []models.Orders
		db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", false)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}
func ws_open_orders(db *gorm.DB) []byte {
	var all_data []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", true)
	data, err := json.Marshal(all_data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

//	func ws_open_orders(db *gorm.DB) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			var all_data []models.Orders
//			db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", false)
//			w.Header().Set("Content-Type", "application/json")
//			json.NewEncoder(w).Encode(all_data)
//		}
//	}
func get_closed_orders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []models.Orders
		db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", true)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}
func update_order(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&data)
		var order models.Orders
		db.Find(&order, "ID=?", data["ID"])
		db.Model(&order).Update("Settled", true)
		db.First(&order, "ID=?", data["ID"])
		fmt.Printf("%+v\n", order)
	}
}
func open_sell_orders(db *gorm.DB) []models.Orders {
	var all_data []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Price in (select price from Orders group by price having count(*)>1) AND Order_Type=? AND Settled=?", "Sell", false)
	// db.Raw("SELECT * FROM models.Orders WHERE Price in (select price from models.Orders group by price having count(*)>1) AND Order_type=?", "Sell").Scan(&all_data)
	return all_data
}

func open_buy_orders(db *gorm.DB) []models.Orders {
	var all_data []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pairg").Find(&all_data, "Price in (select price from Orders group by price having count(*)>1) AND Order_Type=? And Settled=?", "Buy", false)
	return all_data
}

func settle_orders_live(db *gorm.DB) {
	for {
		settle_orders(db)
		time.Sleep(10 * time.Second)

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

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}
func setupwsRoutes(router *mux.Router) {
	pool := websocket.NewPool()
	go pool.Start()
	go pool.Databroadcast()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
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
	// go update_data_live(db)
	go settle_orders_live(db)
	r := mux.NewRouter()
	setupwsRoutes(r)
	r.HandleFunc("/api/all", get_all_data(db)).Methods("GET")
	r.HandleFunc("/api/btc", get_btc_data(db)).Methods("GET")
	r.HandleFunc("/api/eth", get_eth_data(db)).Methods("GET")
	r.HandleFunc("/api/xmr", get_xmr_data(db)).Methods("GET")
	r.HandleFunc("/api/ltc", get_ltc_data(db)).Methods("GET")
	r.HandleFunc("/api/dash", get_dash_data(db)).Methods("GET")
	r.HandleFunc("/api/orders/open", get_open_orders(db)).Methods("GET")
	r.HandleFunc("/api/orders/closed", get_closed_orders(db)).Methods("GET")
	r.HandleFunc("/api/orders/new", new_order(db)).Methods("POST")
	r.HandleFunc("/api/orders/update", update_order(db)).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}
