package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyptocrypto/go_exchange_api/server/models"
	"github.com/hyptocrypto/go_exchange_api/server/utils"
	"github.com/hyptocrypto/go_exchange_api/server/websocket"
	"gorm.io/gorm"
)

func Get_all_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []models.Trading_Pairs
		db.Find(&all_data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}

func Get_btc_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_BTC").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func Get_eth_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_ETH").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func Get_xmr_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_XMR").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func Get_ltc_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_LTC").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func Get_dash_data(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.Trading_Pairs
		db.Where("Ticker=?", "USDT_DASH").First(&data)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
func New_order(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&data)
		amount_val := utils.Interface_to_float(data["Amount"].(string))
		price_val := utils.Interface_to_float(data["Price"].(string))
		settled_val := utils.String_to_bool(data["Settled"].(string))
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
func Get_open_orders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []models.Orders
		db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", false)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}
func Ws_open_orders(db *gorm.DB) []byte {
	var all_data []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", true)
	data, err := json.Marshal(all_data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func Get_closed_orders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var all_data []models.Orders
		db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&all_data, "Settled=?", true)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all_data)
	}
}
func Update_order(db *gorm.DB) http.HandlerFunc {
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
func SetupwsRoutes(router *mux.Router) {
	pool := websocket.NewPool()
	go pool.Start()
	go pool.Databroadcast()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}
