package main

import (
	_ "dcenter/internal/logic"
	"dcenter/utility/gold"

	"github.com/gogf/gf/v2/os/gctx"

	"dcenter/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	_ = gold.NewGoldService()
	cmd.Main.Run(gctx.GetInitCtx())
}
