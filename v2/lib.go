package main

import "github.com/jmoiron/sqlx"

func GetAllOrders(db sqlx.DB) AllOrders {
	var open []OpenOrder
	var closed []ClosedOrder

	open_err := db.Select(&open, "Select * from open_order;")
	closed_err := db.Select(&closed, "Select * from closed_order;")
	if open_err != nil {
		panic(open_err)
	}
	if closed_err != nil {
		panic(closed_err)
	}
	return AllOrders{OpenOrders: open, ClosedOrders: closed}
}
