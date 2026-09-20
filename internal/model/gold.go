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

// GoldPriceHistoryItem 金价历史数据点（每天一条，取当天最后一次报价）。
type GoldPriceHistoryItem struct {
	Date       string  // 日期（2006-01-02）
	Price      float64 // 当日最新价
	OpenPrice  float64 // 开盘价
	HighPrice  float64 // 最高价
	LowPrice   float64 // 最低价
	PrevSettle float64 // 昨结算价
	Unit       string  // 单位
}
