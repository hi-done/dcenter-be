package v1

import "github.com/gogf/gf/v2/frame/g"

type GetGoldPriceReq struct {
	g.Meta `path:"/gold/getGoldPrice" method:"get" summary:"WebSocket 实时推送金价 （每 1 分钟一次"`
}

type GetGoldPriceRes struct {
	Price     float64 `json:"price" dc:"国内金价"`
	Unit      string  `json:"unit" dc:"单位"`
	PriceTime string  `json:"priceTime" dc:"金价时间"`
}
