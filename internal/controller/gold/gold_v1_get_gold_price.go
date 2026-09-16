package gold

import (
	"context"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"

	v1 "dcenter/api/gold/v1"
	"dcenter/internal/service"
)

const (
	// wsPushInterval 金价推送间隔（每 1 分钟一次）
	wsPushInterval = time.Minute
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

// GetGoldPrice 将 HTTP 连接升级为 WebSocket，建立后立即推送一次最新金价，
// 之后每 1 分钟推送一次。推送消息结构与 HTTP 接口保持一致：
// {"code":0,"message":"OK","data":{"price":928.00,"unit":"元/克","priceTime":"2026-09-15 14:20:33"}}。
func (c *ControllerV1) GetGoldPrice(ctx context.Context, req *v1.GetGoldPriceReq) (res *v1.GetGoldPriceRes, err error) {
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

	// 金价由 utility 层按分钟定时采集，接口层负责按固定间隔将最新值推送给前端
	ticker := time.NewTicker(wsPushInterval)
	defer ticker.Stop()
	for {
		price, unit, priceTime := service.Gold().GetGoldPrice()
		_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		if writeErr := conn.WriteJSON(buildPushMessage(price, unit, priceTime)); writeErr != nil {
			g.Log().Debugf(ctx, "gold price websocket closed: %v", writeErr)
			return nil, nil
		}

		select {
		case <-clientCtx.Done():
			// 请求取消或前端断开
			return nil, nil
		case <-ticker.C:
		}
	}
}

// buildPushMessage 将业务层返回的最新金价转换为推送消息。
func buildPushMessage(price float64, unit string, priceTime time.Time) ghttp.DefaultHandlerResponse {
	return ghttp.DefaultHandlerResponse{
		Code:    gcode.CodeOK.Code(),
		Message: gcode.CodeOK.Message(),
		Data: v1.GetGoldPriceRes{
			Price:     price,
			Unit:      unit,
			PriceTime: priceTime.Format("2006-01-02 15:04:05"),
		},
	}
}
