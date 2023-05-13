package models

import (
	"time"

	"gorm.io/gorm"
)

type Trading_Pairs struct {
	gorm.Model
	Ticker         string
	Price          float64
	Daily_Volume   float64
	Daily_High     float64
	Daily_Low      float64
	Percent_Change float64
}

type Orders struct {
	gorm.Model
	CreatedAt       time.Time
	Trading_PairID  uint
	Trading_Pair    Trading_Pairs
	Order_Type      string
	Opening_Amount  float64
	Current_Amount  float64
	Settled         bool
	Partial_Settled bool
	Price           float64
}
