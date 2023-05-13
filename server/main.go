package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyptocrypto/go_exchange_api/server/config"
	"github.com/hyptocrypto/go_exchange_api/server/engine"
	"github.com/hyptocrypto/go_exchange_api/server/handlers"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	go engine.Update_data_live(db)
	go engine.Settle_orders_live(db)
	r := mux.NewRouter()
	handlers.SetupwsRoutes(r)
	r.HandleFunc("/api/all", handlers.Get_all_data(db)).Methods("GET")
	r.HandleFunc("/api/btc", handlers.Get_btc_data(db)).Methods("GET")
	r.HandleFunc("/api/eth", handlers.Get_eth_data(db)).Methods("GET")
	r.HandleFunc("/api/xmr", handlers.Get_xmr_data(db)).Methods("GET")
	r.HandleFunc("/api/ltc", handlers.Get_ltc_data(db)).Methods("GET")
	r.HandleFunc("/api/dash", handlers.Get_dash_data(db)).Methods("GET")
	r.HandleFunc("/api/orders/open", handlers.Get_open_orders(db)).Methods("GET")
	r.HandleFunc("/api/orders/closed", handlers.Get_closed_orders(db)).Methods("GET")
	r.HandleFunc("/api/orders/new", handlers.New_order(db)).Methods("POST")
	r.HandleFunc("/api/orders/update", handlers.Update_order(db)).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}
