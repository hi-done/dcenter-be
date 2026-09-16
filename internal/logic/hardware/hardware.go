package hardware

import (
	"context"
	"math"
	"time"

	"dcenter/internal/model"
	"dcenter/internal/service"
	"dcenter/utility/hardware"
)

type sHardware struct {
}

func init() {
	service.RegisterHardware(New())
}

func New() service.IHardware {
	return &sHardware{}
}

// GetCpuTemp 采样一次 CPU 温度（摄氏度，保留两位小数）。
func (s *sHardware) GetCpuTemp(ctx context.Context) (float64, error) {
	temp, err := hardware.GetCPUTemp()
	if err != nil {
		return 0, err
	}
	return math.Round(temp*100) / 100, nil
}

// WatchCpuTemp 按 interval 周期采样 CPU 温度并写入返回的 channel（建立后立即采样一次）；
// ctx 取消时关闭 channel 并结束采样。
func (s *sHardware) WatchCpuTemp(ctx context.Context, interval time.Duration) <-chan model.CpuTempSample {
	samples := make(chan model.CpuTempSample, 1)
	go func() {
		defer close(samples)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			// 立即采样一次，避免调用方等待首个数据
			sample := s.sample(ctx)
			select {
			case <-ctx.Done():
				return
			case samples <- sample:
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return samples
}

// sample 采样一次并组装为 Sample，供 WatchCpuTemp 周期调用。
func (s *sHardware) sample(ctx context.Context) model.CpuTempSample {
	temp, err := s.GetCpuTemp(ctx)
	if err != nil {
		return model.CpuTempSample{Err: err}
	}
	return model.CpuTempSample{
		TempCelsius: temp,
		Timestamp:   time.Now(),
	}
}
