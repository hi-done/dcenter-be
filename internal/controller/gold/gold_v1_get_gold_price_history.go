package gold

import (
	"context"

	v1 "dcenter/api/gold/v1"
	"dcenter/internal/model"
	"dcenter/internal/service"
)

// GetGoldPriceHistory 返回指定范围（近一个月/半年/一年，默认近一个月）的金价历史数据，
// 数据点按日期升序排列，取自 gold_price 表中每天的最后一次报价。
func (c *ControllerV1) GetGoldPriceHistory(ctx context.Context, req *v1.GetGoldPriceHistoryReq) (res *v1.GetGoldPriceHistoryRes, err error) {
	items, err := service.Gold().GetGoldPriceHistory(ctx, req.Range)
	if err != nil {
		return nil, err
	}
	return &v1.GetGoldPriceHistoryRes{
		Range: req.Range,
		Items: goldPriceHistoryToAPI(items),
	}, nil
}

// goldPriceHistoryToAPI 将业务层历史数据转换为 API 输出结构。
func goldPriceHistoryToAPI(items []model.GoldPriceHistoryItem) []v1.GoldPriceHistoryItem {
	result := make([]v1.GoldPriceHistoryItem, 0, len(items))
	for _, item := range items {
		result = append(result, v1.GoldPriceHistoryItem{
			Date:       item.Date,
			Price:      item.Price,
			OpenPrice:  item.OpenPrice,
			HighPrice:  item.HighPrice,
			LowPrice:   item.LowPrice,
			PrevSettle: item.PrevSettle,
			Unit:       item.Unit,
		})
	}
	return result
}
