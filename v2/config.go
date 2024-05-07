package main

import "github.com/jmoiron/sqlx"

type Conf struct {
	DB_Path string
	Host    string
	Port    string
}

func (c Conf) get_db() sqlx.DB {
	db, err := sqlx.Connect("sqlite3", c.DB_Path)
	if err != nil {
		panic(err)
	}
	return *db
}

var Config = Conf{
	DB_Path: "./exchange.db",
	Host:    "127.0.0.1",
	Port:    ":5000",
}
