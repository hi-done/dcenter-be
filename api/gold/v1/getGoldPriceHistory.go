package v1

import "github.com/gogf/gf/v2/frame/g"

// GetGoldPriceHistoryReq 金价历史数据查询请求。
// range 取值：1m-近一个月（默认）、6m-近半年、1y-近一年。
type GetGoldPriceHistoryReq struct {
	g.Meta `path:"/gold/getGoldPriceHistory" tags:"Gold" method:"get" summary:"获取金价历史数据（近一个月/半年/一年，默认近一个月）"`
	Range  string `json:"range" d:"1m" v:"in:1m,6m,1y" dc:"时间范围：1m-近一个月（默认），6m-近半年，1y-近一年"`
}

// GoldPriceHistoryItem 单日金价历史数据点。
type GoldPriceHistoryItem struct {
	Date       string  `json:"date" dc:"日期（2006-01-02）"`
	Price      float64 `json:"price" dc:"当日最新价"`
	OpenPrice  float64 `json:"openPrice" dc:"开盘价"`
	HighPrice  float64 `json:"highPrice" dc:"最高价"`
	LowPrice   float64 `json:"lowPrice" dc:"最低价"`
	PrevSettle float64 `json:"prevSettle" dc:"昨结算价"`
	Unit       string  `json:"unit" dc:"单位"`
}

// GetGoldPriceHistoryRes 金价历史数据响应。
type GetGoldPriceHistoryRes struct {
	Range string                 `json:"range" dc:"时间范围标识"`
	Items []GoldPriceHistoryItem `json:"items" dc:"金价历史数据点（按时间升序）"`
}
