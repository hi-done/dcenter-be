package model

import "time"

// CpuTempSample CPU 温度采样点。
// 采样成功时携带温度值与采样时间；采样失败时 Err 非空，其余字段无意义。
type CpuTempSample struct {
	TempCelsius float64   // 温度（摄氏度）
	Timestamp   time.Time // 采样时间
	Err         error     // 采样错误
}
