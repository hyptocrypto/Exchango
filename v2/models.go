package main

type TradingPair struct {
	Name  string `db:"name"`
	Price int    `db:"price"`
}

type OpenOrder struct {
	ID               int    `db:"id"`
	TradingPair      string `db:"trading_pair"`
	Amount           int    `db:"amount"`
	Price            int    `db:"price"`
	Type             string `db:"type"`
	PartiallySettled bool   `db:"partially_settled"`
}

type ClosedOrder struct {
	ID          int `db:"id"`
	OpenOrderID int `db:"open_order"`
}

type AllOrders struct {
	OpenOrders   []OpenOrder
	ClosedOrders []ClosedOrder
}
