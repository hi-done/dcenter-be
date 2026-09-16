// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GoldPrice is the golang structure for table gold_price.
type GoldPrice struct {
	Id         int64       `json:"id"         orm:"id"          description:""` //
	Symbol     string      `json:"symbol"     orm:"symbol"      description:""` //
	Name       string      `json:"name"       orm:"name"        description:""` //
	Price      float64     `json:"price"      orm:"price"       description:""` //
	OpenPrice  float64     `json:"openPrice"  orm:"open_price"  description:""` //
	HighPrice  float64     `json:"highPrice"  orm:"high_price"  description:""` //
	LowPrice   float64     `json:"lowPrice"   orm:"low_price"   description:""` //
	PrevSettle float64     `json:"prevSettle" orm:"prev_settle" description:""` //
	Unit       string      `json:"unit"       orm:"unit"        description:""` //
	QuoteTime  *gtime.Time `json:"quoteTime"  orm:"quote_time"  description:""` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""` //
}
