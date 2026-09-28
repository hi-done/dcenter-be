package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"dcenter/internal/dao"
	"dcenter/internal/model/do"
	"dcenter/internal/model/entity"
	"dcenter/utility/coding"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SignIn 用户登录：校验用户名密码与账号状态，签发 access/refresh 双 token，返回新 token 对
func (s *sUser) SignIn(ctx context.Context, userName, password string) (string, string, error) {
	// 1. 查询用户（软删除自动过滤）
	var user entity.DcUser
	err := dao.DcUser.Ctx(ctx).Where(do.DcUser{UserName: userName}).Scan(&user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", "", gerror.Wrap(err, "查询用户失败")
	}
	// 用户不存在与密码错误返回同一文案，避免暴露账号是否存在
	if user.Id == 0 {
		return "", "", gerror.NewCode(gcode.CodeValidationFailed, "用户名或密码错误")
	}

	// 2. 密码校验
	if !coding.VerifyPassword(password, user.Password) {
		return "", "", gerror.NewCode(gcode.CodeValidationFailed, "用户名或密码错误")
	}

	// 3. 账号状态校验：1正常 0禁用 2锁定
	if user.Status != 1 {
		msg := "账号状态异常，无法登录"
		if user.Status == 0 {
			msg = "账号已被禁用"
		} else if user.Status == 2 {
			msg = "账号已被锁定"
		}
		return "", "", gerror.NewCode(gcode.CodeValidationFailed, msg)
	}

	// 4. 签发 access/refresh 双 token
	accessToken, err := coding.GenerateAccessToken(int64(user.Id), user.UserName)
	if err != nil {
		return "", "", gerror.Wrap(err, "签发 Access Token 失败")
	}
	refreshToken, err := coding.GenerateRefreshToken(int64(user.Id), user.UserName)
	if err != nil {
		return "", "", gerror.Wrap(err, "签发 Refresh Token 失败")
	}

	// 5. Refresh Token 落 Redis（覆盖旧会话），TTL 与有效期一致，供刷新轮换与吊销校验
	key := fmt.Sprintf("refresh_token:%d", user.Id)
	ttlSeconds := int64(coding.RefreshTokenExpire(ctx) / time.Second)
	if err = g.Redis().SetEX(ctx, key, refreshToken, ttlSeconds); err != nil {
		return "", "", gerror.Wrap(err, "保存登录状态失败")
	}

	// 6. 记录最近登录时间与 IP，失败不影响登录结果
	if _, err = dao.DcUser.Ctx(ctx).Where(do.DcUser{Id: user.Id}).Data(do.DcUser{
		LastLoginAt: gtime.Now(),
		LastLoginIp: clientIp(ctx),
	}).Update(); err != nil {
		g.Log().Warningf(ctx, "更新最近登录信息失败, userId=%d: %v", user.Id, err)
	}

	return accessToken, refreshToken, nil
}

// clientIp 从请求上下文提取客户端 IP，非 HTTP 上下文返回空串
func clientIp(ctx context.Context) string {
	if r := g.RequestFromCtx(ctx); r != nil {
		return r.GetClientIp()
	}
	return ""
}
