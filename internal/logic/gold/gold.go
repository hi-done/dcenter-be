package gold

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"dcenter/internal/dao"
	"dcenter/internal/model"
	"dcenter/internal/service"
	"dcenter/utility/gold"
)

var goldService = gold.NewGoldService()

// 金价历史查询范围标识。
const (
	historyRangeMonth    = "1m" // 近一个月（默认）
	historyRangeHalfYear = "6m" // 近半年
	historyRangeYear     = "1y" // 近一年
)

// goldHistorySymbol 历史查询的品种：国内黄金延期（与实时金价接口保持一致）。
const goldHistorySymbol = "AUTD"

type sGold struct {
}

func init() {
	service.RegisterGold(New())
}

func New() service.IGold {
	return &sGold{}
}

func (s *sGold) GetGoldPrice() (float64, string, time.Time) {
	return goldService.GetGoldPrice()
}

// GetGoldPriceHistory 获取金价历史数据 近一个月，半年，一年 默认近一个月。
// period 取值 1m/6m/1y，返回按日期升序排列的数据点（每天一条）。
func (s *sGold) GetGoldPriceHistory(ctx context.Context, period string) ([]model.GoldPriceHistoryItem, error) {
	startTime, err := historyStartTime(period)
	if err != nil {
		return nil, err
	}

	columns := dao.GoldPrice.Columns()
	records, err := dao.GoldPrice.Ctx(ctx).
		Where(columns.Symbol, goldHistorySymbol).
		WhereGTE(columns.QuoteTime, startTime).
		OrderAsc(columns.QuoteTime).
		All()
	if err != nil {
		return nil, err
	}

	items := make([]model.GoldPriceHistoryItem, 0, len(records))
	for _, record := range records {
		item := model.GoldPriceHistoryItem{
			Price:      record[columns.Price].Float64(),
			OpenPrice:  record[columns.OpenPrice].Float64(),
			HighPrice:  record[columns.HighPrice].Float64(),
			LowPrice:   record[columns.LowPrice].Float64(),
			PrevSettle: record[columns.PrevSettle].Float64(),
			Unit:       record[columns.Unit].String(),
		}
		if quoteTime := record[columns.QuoteTime].GTime(); quoteTime != nil {
			// 注意：gtime.Time 的 Format 是 PHP 风格占位符（如 Y-m-d），标准库 layout 需用 Layout。
			item.Date = quoteTime.Layout("2006-01-02")
		}
		items = append(items, item)
	}
	return items, nil
}

// historyStartTime 将查询范围标识转换为起始时间（从起始日期的 00:00 开始），
// 空标识按近一个月处理，未知标识返回参数错误。
func historyStartTime(period string) (time.Time, error) {
	now := time.Now()
	var start time.Time
	switch period {
	case "", historyRangeMonth:
		start = now.AddDate(0, -1, 0)
	case historyRangeHalfYear:
		start = now.AddDate(0, -6, 0)
	case historyRangeYear:
		start = now.AddDate(-1, 0, 0)
	default:
		return time.Time{}, gerror.NewCodef(gcode.CodeInvalidParameter,
			"invalid history range: %q, expect one of 1m/6m/1y", period)
	}
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location()), nil
}
