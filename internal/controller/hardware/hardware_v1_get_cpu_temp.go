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
	"dcenter/internal/model"
	"dcenter/internal/service"
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

	// 业务层按固定间隔采样 CPU 温度，接口层负责将采样结果推送给前端
	samples := service.Hardware().WatchCpuTemp(clientCtx, wsPushInterval)
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case sample, ok := <-samples:
			if !ok {
				return nil, nil
			}
			_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if writeErr := conn.WriteJSON(buildPushMessage(sample)); writeErr != nil {
				g.Log().Debugf(ctx, "cpu temp websocket closed: %v", writeErr)
				return nil, nil
			}
		}
	}
}

// buildPushMessage 将业务层采样结果转换为推送消息；
// 采样失败时返回带错误码的消息，前端可据此展示异常状态。
func buildPushMessage(sample model.CpuTempSample) ghttp.DefaultHandlerResponse {
	if sample.Err != nil {
		return ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: sample.Err.Error(),
		}
	}
	return ghttp.DefaultHandlerResponse{
		Code:    gcode.CodeOK.Code(),
		Message: gcode.CodeOK.Message(),
		Data: v1.GetCpuTempRes{
			TempCelsius: sample.TempCelsius,
			Timestamp:   sample.Timestamp.UnixMilli(),
		},
	}
}
