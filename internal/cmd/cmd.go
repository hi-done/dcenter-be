package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"dcenter/internal/controller/gold"
	"dcenter/internal/controller/hardware"
	"dcenter/internal/controller/sys"
	"dcenter/internal/controller/user"
	"dcenter/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Use(service.Middleware().CORS)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hardware.NewV1(),
					gold.NewV1(),
					sys.NewV1(),
				)
				// 需要登录认证的路由
				group.Middleware(service.Middleware().Auth)
				group.Bind(
					user.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
