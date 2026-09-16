// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GoldPrice is the golang structure of table gold_price for DAO operations like Where/Data.
type GoldPrice struct {
	g.Meta     `orm:"table:gold_price, do:true"`
	Id         any         //
	Symbol     any         //
	Name       any         //
	Price      any         //
	OpenPrice  any         //
	HighPrice  any         //
	LowPrice   any         //
	PrevSettle any         //
	Unit       any         //
	QuoteTime  *gtime.Time //
	CreatedAt  *gtime.Time //
}
