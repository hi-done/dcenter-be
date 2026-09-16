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
	IHardware interface {
		// GetCpuTemp 采样一次 CPU 温度（摄氏度，保留两位小数）。
		GetCpuTemp(ctx context.Context) (float64, error)
		// WatchCpuTemp 按 interval 周期采样 CPU 温度并写入返回的 channel（建立后立即采样一次）；
		// ctx 取消时关闭 channel 并结束采样。
		WatchCpuTemp(ctx context.Context, interval time.Duration) <-chan model.CpuTempSample
	}
)

var (
	localHardware IHardware
)

func Hardware() IHardware {
	if localHardware == nil {
		panic("implement not found for interface IHardware, forgot register?")
	}
	return localHardware
}

func RegisterHardware(i IHardware) {
	localHardware = i
}
