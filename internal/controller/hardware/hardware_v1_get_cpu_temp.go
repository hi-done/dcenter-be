package hardware

import (
	"context"
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

const (
	// wsPushInterval CPU 温度推送间隔（每 1 秒一次）
	wsPushInterval = time.Second
	// wsWriteWait 单次写入超时时间
	wsWriteWait = 5 * time.Second
)

// wsUpGrader WebSocket 升级器。
// 本地看板场景下放开来源校验；握手失败保持静默，由业务层统一返回 JSON 错误结构。
var wsUpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
	},
}

// GetCpuTemp 将 HTTP 连接升级为 WebSocket，并每秒向前端推送一次系统 CPU 温度。
// 温度直接取自全局唯一的采样服务（utility/hardware.NewHardwareService）：
// 所有连接共享同一份采样结果，连接结束后自动释放引用，无调用方时采样自动停止。
// 推送消息结构与 HTTP 接口保持一致：{"code":0,"message":"OK","data":{"tempCelsius":27.85,"timestamp":1757836800000}}。
func (c *ControllerV1) GetCpuTemp(ctx context.Context, req *v1.GetCpuTempReq) (res *v1.GetCpuTempRes, err error) {
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
	if pushErr := pushCpuTemp(ctx, conn, hs); pushErr != nil {
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
			if pushErr := pushCpuTemp(ctx, conn, hs); pushErr != nil {
				return nil, nil
			}
		}
	}
}

// pushCpuTemp 读取一次最新采样结果并写入连接；写入失败时返回错误。
func pushCpuTemp(ctx context.Context, conn *websocket.Conn, hs *hardware.HardwareService) error {
	tempCelsius, sampleErr := hs.GetInfo()
	_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
	if err := conn.WriteJSON(buildPushMessage(tempCelsius, sampleErr)); err != nil {
		g.Log().Debugf(ctx, "cpu temp websocket closed: %v", err)
		return err
	}
	return nil
}

// buildPushMessage 将一次采样结果转换为推送消息；
// 采样失败（或首个采样尚未完成）时返回带错误码的消息，前端可据此展示异常状态。
func buildPushMessage(tempCelsius float64, sampleErr error) ghttp.DefaultHandlerResponse {
	if sampleErr != nil {
		return ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: sampleErr.Error(),
		}
	}
	return ghttp.DefaultHandlerResponse{
		Code:    gcode.CodeOK.Code(),
		Message: gcode.CodeOK.Message(),
		Data: v1.GetCpuTempRes{
			TempCelsius: tempCelsius,
			Timestamp:   time.Now().UnixMilli(),
		},
	}
}
