package hardware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"

	v1 "dcenter/api/hardware/v1"
	"dcenter/utility/hardware"
)

// GetSysInfo 将 HTTP 连接升级为 WebSocket，并每秒向前端推送一次系统信息
// （CPU 温度、开机时长、内存与系统盘使用情况）。
// 数据直接取自全局唯一的采样服务（utility/hardware.NewHardwareService）：
// 所有连接共享同一份采样结果，连接结束后自动释放引用，无调用方时采样自动停止。
// 推送消息结构：{"code":0,"message":"OK","data":{"tempCelsius":27.85,"uptimeSeconds":234567,
// "memory":{"totalBytes":34127908864,"usedBytes":14765535232,"usedPercent":43.27},
// "disk":{...},"timestamp":1757836800000}}。
func (c *ControllerV1) GetSysInfo(ctx context.Context, req *v1.GetSysInfoReq) (res *v1.GetSysInfoRes, err error) {
	r := g.RequestFromCtx(ctx)
	conn, err := wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		// 非 WebSocket 握手请求（例如浏览器直接访问），返回 426 与统一的 JSON 错误结构
		r.Response.WriteStatus(http.StatusUpgradeRequired)
		return nil, gerror.NewCode(gcode.CodeInvalidRequest, "websocket connection required: "+err.Error())
	}
	defer conn.Close()

	// WebSocket 协议要求持续读取连接才能处理 Close/Ping/Pong 控制帧，
	// 同时借此感知前端断开并结束推送。
	clientCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer cancel()
		for {
			if _, _, readErr := conn.ReadMessage(); readErr != nil {
				return
			}
		}
	}()

	// 获取全局唯一的硬件采样服务；clientCtx 结束时自动释放本次引用
	hs := hardware.NewHardwareService(clientCtx)

	// 建立连接后立即推送一次，避免前端等待首个推送周期
	if pushErr := pushSysInfo(ctx, conn, hs); pushErr != nil {
		return nil, nil
	}

	// 按固定间隔读取最新采样结果并推送给前端
	ticker := time.NewTicker(wsPushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-clientCtx.Done():
			return nil, nil
		case <-ticker.C:
			if pushErr := pushSysInfo(ctx, conn, hs); pushErr != nil {
				return nil, nil
			}
		}
	}
}

// pushSysInfo 读取一次最新采样结果并写入连接；写入失败时返回错误。
func pushSysInfo(ctx context.Context, conn *websocket.Conn, hs *hardware.HardwareService) error {
	tempCelsius, tempErr := hs.GetInfo()
	sample, sampleErr := hs.GetSysInfo()
	_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
	if err := conn.WriteJSON(buildSysInfoPushMessage(tempCelsius, tempErr, sample, sampleErr)); err != nil {
		g.Log().Debugf(ctx, "sys info websocket closed: %v", err)
		return err
	}
	return nil
}

// buildSysInfoPushMessage 将一次采样结果转换为推送消息；
// 任一项采样失败（或首个采样尚未完成）时返回带错误码的消息，前端可据此展示异常状态。
func buildSysInfoPushMessage(
	tempCelsius float64, tempErr error,
	sample hardware.SysInfoSample, sampleErr error,
) ghttp.DefaultHandlerResponse {
	if err := errors.Join(tempErr, sampleErr, sample.UptimeErr, sample.MemoryErr, sample.DiskErr); err != nil {
		return ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: err.Error(),
		}
	}
	return ghttp.DefaultHandlerResponse{
		Code:    gcode.CodeOK.Code(),
		Message: gcode.CodeOK.Message(),
		Data: v1.GetSysInfoRes{
			TempCelsius:   tempCelsius,
			UptimeSeconds: int64(sample.Uptime / time.Second),
			Memory:        usageStatToAPI(sample.Memory),
			Disk:          usageStatToAPI(sample.Disk),
			Timestamp:     time.Now().UnixMilli(),
		},
	}
}

// usageStatToAPI 将底层采样结果转换为 API 输出结构。
func usageStatToAPI(stat hardware.UsageStat) v1.UsageStat {
	return v1.UsageStat{
		TotalBytes:  stat.Total,
		UsedBytes:   stat.Used,
		UsedPercent: stat.UsedPercent,
	}
}
