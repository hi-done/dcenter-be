package user

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// SignOut 用户注销：删除 Redis 中的 refresh token，使当前会话失效
func (s *sUser) SignOut(ctx context.Context) error {
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		return gerror.NewCode(gcode.CodeNotAuthorized, "用户未登录")
	}
	// 删除会话 refresh token（不存在也视为成功，幂等）
	key := fmt.Sprintf("refresh_token:%d", userId)
	if _, err := g.Redis().Del(ctx, key); err != nil {
		return gerror.Wrap(err, "注销失败")
	}
	return nil
}
