// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"dcenter/internal/model"
	"time"
)

type (
	IGold interface {
		GetGoldPrice() (float64, string, time.Time)
		// GetGoldPriceHistory 获取金价历史数据 近一个月，半年，一年 默认近一个月。
		// period 取值 1m/6m/1y，返回按日期升序排列的数据点（每天一条）。
		GetGoldPriceHistory(ctx context.Context, period string) ([]model.GoldPriceHistoryItem, error)
	}
)

var (
	localGold IGold
)

func Gold() IGold {
	if localGold == nil {
		panic("implement not found for interface IGold, forgot register?")
	}
	return localGold
}

func RegisterGold(i IGold) {
	localGold = i
}
