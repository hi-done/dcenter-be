package middleware

import (
	"dcenter/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

type sMiddleware struct{}

func init() {
	service.RegisterMiddleware(New())
}

func New() *sMiddleware {
	return &sMiddleware{}
}

func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	// 双 token 刷新机制所需的自定义头：
	//   X-Refresh-Token 随请求携带，Access Token 过期时由 Auth 中间件消费；
	//   轮换后的新 token 经响应头 X-Access-Token / X-Refresh-Token 下发，
	//   需通过 Access-Control-Expose-Headers 暴露给跨域前端读取。
	r.Response.Header().Set("Access-Control-Allow-Headers",
		"Origin,Content-Type,Accept,Authorization,X-Refresh-Token")
	r.Response.Header().Set("Access-Control-Expose-Headers", "X-Access-Token,X-Refresh-Token")
	r.Middleware.Next()
}
