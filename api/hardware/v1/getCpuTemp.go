package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetCpuTempReq WebSocket 握手请求。
// 前端使用 ws://<host>/hardware/cpu-temp 建立连接后，服务端每秒推送一次 CPU 温度。
type GetCpuTempReq struct {
	g.Meta `path:"/hardware/cpu-temp" method:"get" tags:"Hardware" summary:"WebSocket 实时推送系统 CPU 温度（每 1 秒一次）"`
}

// GetCpuTempRes WebSocket 推送的 CPU 温度数据。
type GetCpuTempRes struct {
	TempCelsius float64 `json:"tempCelsius" dc:"CPU 温度（摄氏度）"`
	Timestamp   int64   `json:"timestamp" dc:"采样时间戳（毫秒）"`
}
