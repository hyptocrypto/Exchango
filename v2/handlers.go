package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db sqlx.DB
var ws_pool *Pool
var update_channel = make(chan bool)

var ws_upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Init db and start websocket pool
func init() {
	db = Config.get_db()
	ws_pool = NewPool(&db)
	go ws_pool.Start(&db, update_channel)
}

func GetPairs(w http.ResponseWriter, r *http.Request) {
	conn, err := ws_upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-update_channel:
			var pairs []TradingPair
			err := db.Select(&pairs, "Select name, price from trading_pair;")
			if err != nil {
				log.Fatal(err)
			}
			ret, err := json.Marshal(pairs)
			if err != nil {
				log.Fatal(err)
			}
			conn.WriteMessage(websocket.TextMessage, ret)
		}
	}

}

func Update(w http.ResponseWriter, r *http.Request) {
	update_channel <- true
}

func WsOpenOrders(w http.ResponseWriter, r *http.Request) {
	conn, err := ws_upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := NewClient(conn, ws_pool)
	client.Pool.Register <- client
	update_channel <- true
	// client.PushMsg(GetAllOrders(*client.Pool.DB))
}

func WsClosedOrders(w http.ResponseWriter, r *http.Request) {
	var closed_orders []ClosedOrder
	err := db.Select(&closed_orders, "Select * from closed_order;")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, closed_orders)
}
