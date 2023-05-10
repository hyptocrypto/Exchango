package websocket

import (
	"fmt"
	"time"

	"github.com/hyptocrypto/go_exchange_api/server/config"
	"github.com/hyptocrypto/go_exchange_api/server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	var open_orders []models.Orders
	var closed_orders []models.Orders
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&open_orders, "Settled=?", false)
	db.Model(&models.Orders{}).Preload("Trading_Pair").Find(&closed_orders, "Settled=?", true)
	all_orders := map[string]interface{}{"open_orders": open_orders, "closed_orders": closed_orders}
	return all_orders
}

func (pool *Pool) Databroadcast() {
	db, err := gorm.Open(sqlite.Open(config.DB_Path), &gorm.Config{})
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
