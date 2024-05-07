package main

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

type Message struct {
	Type int    `json:"type"`
	Body []byte `json:"body"`
}

type Client struct {
	ID   uuid.UUID
	Conn *websocket.Conn
	Pool *Pool
}

func NewClient(conn *websocket.Conn, pool *Pool) *Client {
	return &Client{
		ID:   uuid.New(),
		Conn: conn,
		Pool: pool,
	}
}

func (c *Client) PushMsg(msg interface{}) {
	err := c.Conn.WriteJSON(msg)
	if err != nil {
		c.Pool.Unregister <- c
	}
}

type Pool struct {
	DB         *sqlx.DB
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mutex      sync.RWMutex
}

func NewPool(db *sqlx.DB) *Pool {
	return &Pool{
		DB:         db,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		mutex:      sync.RWMutex{},
	}
}

// Push message to all clients in pool
func (p *Pool) PushMsgs(msg interface{}) {
	for client := range p.Clients {
		client.PushMsg(msg)
	}
}

func (pool *Pool) Start(db *sqlx.DB, update_channel chan bool) {
	// Start a goroutine per channel.
	// This is done instead of creating buffered channels.
	// Otherwise we can run into deadlocks where we write to a channel and nothing is reading from it.
	go func() {
		for {
			<-update_channel
			orders := GetAllOrders(*pool.DB)
			pool.PushMsgs(orders)
		}
	}()
	go func() {
		for {
			client := <-pool.Register
			pool.mutex.Lock()
			pool.Clients[client] = true
			pool.mutex.Unlock()
			log.Println("New Client: ", client.ID)
			log.Println("Size of Connection Pool: ", len(pool.Clients))
		}
	}()
	go func() {
		for {
			client := <-pool.Unregister
			log.Printf("Purging client: %v from pool.", client.ID)
			pool.mutex.Lock()
			delete(pool.Clients, client)
			pool.mutex.Unlock()
			client.Conn.Close()
		}
	}()
}

// func (pool *Pool) Start(db *sqlx.DB, update_channel chan bool) {
// 	for {
// 		select {
// 		// Push orders to all clients
// 		case <-update_channel:
// 			orders := GetAllOrders(*pool.DB)
// 			ws_pool.PushMsgs(orders)

// 		case client := <-pool.Register:
// 			pool.Clients[client] = true
// 			log.Println("New Client: ", client.ID)
// 			log.Println("Size of Connection Pool: ", len(pool.Clients))

// 		case client := <-pool.Unregister:
// 			log.Printf("Purging client: %v from pool.", client.ID)
// 			delete(client.Pool.Clients, client)
// 			client.Conn.Close()

// 		default:
// 			continue
// 		}
// 	}
// }
