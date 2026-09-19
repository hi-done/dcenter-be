package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetSysInfoReq WebSocket 握手请求。
// 前端使用 ws://<host>/hardware/sys-info 建立连接后，服务端每秒推送一次系统信息。
type GetSysInfoReq struct {
	g.Meta `path:"/hardware/sys-info" method:"get" tags:"Hardware" summary:"WebSocket 实时推送系统信息（每 1 秒一次）"`
}

// UsageStat 资源使用量（内存 / 磁盘共用）。
type UsageStat struct {
	TotalBytes  uint64  `json:"totalBytes" dc:"总大小（字节）"`
	UsedBytes   uint64  `json:"usedBytes" dc:"已使用大小（字节）"`
	UsedPercent float64 `json:"usedPercent" dc:"使用率（百分比，保留两位小数）"`
}

// GetSysInfoRes WebSocket 推送的系统信息数据（含 CPU 温度）。
type GetSysInfoRes struct {
	TempCelsius   float64   `json:"tempCelsius" dc:"CPU 温度（摄氏度）"`
	UptimeSeconds int64     `json:"uptimeSeconds" dc:"系统开机时长（秒）"`
	Memory        UsageStat `json:"memory" dc:"内存使用情况"`
	Disk          UsageStat `json:"disk" dc:"系统盘使用情况"`
	Timestamp     int64     `json:"timestamp" dc:"推送时间戳（毫秒）"`
}
