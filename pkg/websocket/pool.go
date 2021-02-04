package websocket

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
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
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// 	fmt.Println(client)
			// 	client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			// }
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break

		}
	}
}

func ws_open_orders(db *gorm.DB) map[string]interface{} {
	var open_orders []orders
	var closed_orders []orders
	db.Model(&orders{}).Preload("Trading_Pair").Find(&open_orders, "Settled=?", false)
	db.Model(&orders{}).Preload("Trading_Pair").Find(&closed_orders, "Settled=?", true)
	all_orders := map[string]interface{}{"open_orders": open_orders, "closed_orders": closed_orders}
	return all_orders
}

func (pool *Pool) Databroadcast() {
	db, err := gorm.Open(sqlite.Open("../mock_exchange.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, sqlite_err := db.DB()
	if sqlite_err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	for t := range time.Tick(5 * time.Second) {
		data := ws_open_orders(db)
		for client, _ := range pool.Clients {
			if err := client.Conn.WriteJSON(data); err != nil {
				client.Pool.Unregister <- client
				delete(pool.Clients, client)
				fmt.Println("Client Pruged at: ", t)
			}
		}
	}
}
