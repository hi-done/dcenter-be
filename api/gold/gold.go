// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package gold

import (
	"context"

	"dcenter/api/gold/v1"
)

type IGoldV1 interface {
	GetGoldPrice(ctx context.Context, req *v1.GetGoldPriceReq) (res *v1.GetGoldPriceRes, err error)
	GetGoldPriceHistory(ctx context.Context, req *v1.GetGoldPriceHistoryReq) (res *v1.GetGoldPriceHistoryRes, err error)
}
