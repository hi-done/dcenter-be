package model

import "time"

// GoldPrice represents a row to be inserted into the gold_price table.
type GoldPrice struct {
	Symbol     string
	Name       string
	Price      float64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	PrevSettle float64
	Unit       string
	QuoteTime  time.Time
}
