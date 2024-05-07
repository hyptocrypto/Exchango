package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(db sqlx.DB) {
	pair_schema := `CREATE TABLE IF NOT EXISTS trading_pair (
        name text primary key,
        price integer
    );`

	open_order_schema := `CREATE TABLE IF NOT EXISTS open_order (
        id integer primary key AUTOINCREMENT,
        trading_pair text,
        amount integer,
		price integer,
        type text check(type in ('buy', 'sell')),
        partially_settled boolean DEFAULT false,
        FOREIGN KEY(trading_pair) REFERENCES trading_pair(name)
    );`

	closed_order_schema := `CREATE TABLE IF NOT EXISTS closed_order (
        id integer primary key AUTOINCREMENT,
		open_order integer,
        FOREIGN KEY(open_order) REFERENCES open_order(id)
    );`
	open_order := `INSERT INTO open_order (trading_pair, amount, price, type) values ('BTC/USD', 1000, 65000, 'buy');`

	db.MustExec(pair_schema)
	db.MustExec(open_order_schema)
	db.MustExec(closed_order_schema)

	tickers := []string{"BTC/USD", "ETH/USD", "XRP/USD", "XMR/USD"}
	for _, ticker := range tickers {
		price := 0 // Replace this with the actual price if you have it
		_, err := db.Exec(`INSERT OR IGNORE INTO trading_pair (name, price) VALUES (?, ?)`, ticker, price)
		if err != nil {
			panic(err)
		}
	}
	db.MustExec(open_order)
}
