package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Trading_Pair struct {
	gorm.Model
	Ticker         string
	Price          float64
	Daily_Volume   float64
	Daily_High     float64
	Daily_Low      float64
	Precent_Change float64
}

type Order struct {
	gorm.Model
	Trading_PairID uint
	Trading_Pair   Trading_Pair
	Order_Type     string
	Amount         float64
	Settled        bool
}

func update_data_live(db *gorm.DB) {
	for {
		update_data(db)
		time.Sleep(60 * time.Second)
	}
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
	var btc = data["USDT_BTC"]
	var eth = data["USDT_ETH"]
	var xmr = data["USDT_XMR"]
	var ltc = data["USDT_LTC"]
	var dash = data["USDT_DASH"]

	db.Model(&Trading_Pair{}).Where("Ticker = ?", "USDT_BTC").Updates(Trading_Pair{
		Price:          interface_to_float(btc.(map[string]interface{})["last"]),
		Daily_Volume:   interface_to_float(btc.(map[string]interface{})["baseVolume"]),
		Daily_High:     interface_to_float(btc.(map[string]interface{})["high24hr"]),
		Daily_Low:      interface_to_float(btc.(map[string]interface{})["low24hr"]),
		Precent_Change: interface_to_float(btc.(map[string]interface{})["percentChange"]),
	})
	db.Model(&Trading_Pair{}).Where("Ticker = ?", "USDT_ETH").Updates(Trading_Pair{
		Price:          interface_to_float(eth.(map[string]interface{})["last"]),
		Daily_Volume:   interface_to_float(eth.(map[string]interface{})["baseVolume"]),
		Daily_High:     interface_to_float(eth.(map[string]interface{})["high24hr"]),
		Daily_Low:      interface_to_float(eth.(map[string]interface{})["low24hr"]),
		Precent_Change: interface_to_float(eth.(map[string]interface{})["percentChange"]),
	})
	db.Model(&Trading_Pair{}).Where("Ticker = ?", "USDT_XMR").Updates(Trading_Pair{
		Price:          interface_to_float(xmr.(map[string]interface{})["last"]),
		Daily_Volume:   interface_to_float(xmr.(map[string]interface{})["baseVolume"]),
		Daily_High:     interface_to_float(xmr.(map[string]interface{})["high24hr"]),
		Daily_Low:      interface_to_float(xmr.(map[string]interface{})["low24hr"]),
		Precent_Change: interface_to_float(xmr.(map[string]interface{})["percentChange"]),
	})
	db.Model(&Trading_Pair{}).Where("Ticker = ?", "USDT_LTC").Updates(Trading_Pair{
		Price:          interface_to_float(ltc.(map[string]interface{})["last"]),
		Daily_Volume:   interface_to_float(ltc.(map[string]interface{})["baseVolume"]),
		Daily_High:     interface_to_float(ltc.(map[string]interface{})["high24hr"]),
		Daily_Low:      interface_to_float(ltc.(map[string]interface{})["low24hr"]),
		Precent_Change: interface_to_float(ltc.(map[string]interface{})["percentChange"]),
	})
	db.Model(&Trading_Pair{}).Where("Ticker = ?", "USDT_DASH").Updates(Trading_Pair{
		Price:          interface_to_float(dash.(map[string]interface{})["last"]),
		Daily_Volume:   interface_to_float(dash.(map[string]interface{})["baseVolume"]),
		Daily_High:     interface_to_float(dash.(map[string]interface{})["high24hr"]),
		Daily_Low:      interface_to_float(dash.(map[string]interface{})["low24hr"]),
		Precent_Change: interface_to_float(dash.(map[string]interface{})["percentChange"]),
	})

}

func interface_to_float(data interface{}) float64 {
	f, err := strconv.ParseFloat(data.(string), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func get_all_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []Trading_Pair
		db.Find(&all_data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}

func get_btc_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []Trading_Pair
		db.Where("Ticker=?", "USDT_BTC").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func get_eth_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []Trading_Pair
		db.Where("Ticker=?", "USDT_ETH").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func get_xmr_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []Trading_Pair
		db.Where("Ticker=?", "USDT_XMR").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func get_ltc_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []Trading_Pair
		db.Where("Ticker=?", "USDT_LTC").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func get_dash_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []Trading_Pair
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
		var trading_pair Trading_Pair
		db.First(&trading_pair, "Ticker=?", data["Trading_Pair"])
		order := Order{Trading_PairID: trading_pair.ID,
			Trading_Pair: trading_pair,
			Order_Type:   data["Order_type"].(string),
			Amount:       data["Amount"].(float64),
			Settled:      data["Settled"].(bool)}
		db.Create(&order)

	}
}
func get_all_orders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []Order
		db.Model(&Order{}).Preload("Trading_Pair").Find(&all_data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}
func update_order(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&data)
		var order Order
		db.Find(&order, "ID=?", data["ID"])
		db.Model(&order).Update("Settled", true)
		db.First(&order, "ID=?", data["ID"])
		fmt.Printf("%+v\n", order)
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("../mock_exchange.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, sqlite_err := db.DB()
	if sqlite_err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	go update_data_live(db)
	r := mux.NewRouter()
	r.HandleFunc("/api/all", get_all_data(db)).Methods("GET")
	r.HandleFunc("/api/btc", get_btc_data(db)).Methods("GET")
	r.HandleFunc("/api/eth", get_eth_data(db)).Methods("GET")
	r.HandleFunc("/api/xmr", get_xmr_data(db)).Methods("GET")
	r.HandleFunc("/api/ltc", get_ltc_data(db)).Methods("GET")
	r.HandleFunc("/api/dash", get_dash_data(db)).Methods("GET")
	r.HandleFunc("/api/orders/all", get_all_orders(db)).Methods("GET")
	r.HandleFunc("/api/orders/new", new_order(db)).Methods("POST")
	r.HandleFunc("/api/orders/update", update_order(db)).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", r))

}
