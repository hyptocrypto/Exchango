package main

import (
	"log"
	"net/http"

	"github.com/hyptocrypto/exchangego/util"
)

func main() {
	util.InitDB(Config.get_db())
	log.Println("Starting Server")
	http.HandleFunc("/pairs", GetPairs)
	http.HandleFunc("/ws/open-orders", WsOpenOrders)
	http.HandleFunc("/ws/closed-orders", WsClosedOrders)
	http.HandleFunc("/update", Update)
	http.ListenAndServe(Config.Host+Config.Port, nil)
}
