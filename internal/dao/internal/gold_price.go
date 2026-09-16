// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GoldPriceDao is the data access object for the table gold_price.
type GoldPriceDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  GoldPriceColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// GoldPriceColumns defines and stores column names for the table gold_price.
type GoldPriceColumns struct {
	Id         string //
	Symbol     string //
	Name       string //
	Price      string //
	OpenPrice  string //
	HighPrice  string //
	LowPrice   string //
	PrevSettle string //
	Unit       string //
	QuoteTime  string //
	CreatedAt  string //
}

// goldPriceColumns holds the columns for the table gold_price.
var goldPriceColumns = GoldPriceColumns{
	Id:         "id",
	Symbol:     "symbol",
	Name:       "name",
	Price:      "price",
	OpenPrice:  "open_price",
	HighPrice:  "high_price",
	LowPrice:   "low_price",
	PrevSettle: "prev_settle",
	Unit:       "unit",
	QuoteTime:  "quote_time",
	CreatedAt:  "created_at",
}

// NewGoldPriceDao creates and returns a new DAO object for table data access.
func NewGoldPriceDao(handlers ...gdb.ModelHandler) *GoldPriceDao {
	return &GoldPriceDao{
		group:    "default",
		table:    "gold_price",
		columns:  goldPriceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GoldPriceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GoldPriceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GoldPriceDao) Columns() GoldPriceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GoldPriceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GoldPriceDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *GoldPriceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
