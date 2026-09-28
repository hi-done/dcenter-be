package user

import (
	"context"
	"fmt"

	"dcenter/internal/dao"
	"dcenter/internal/model/do"
	"dcenter/utility/coding"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// SignUp 用户注册：校验邀请码、加密密码并创建账号，返回新用户ID
func (s *sUser) SignUp(ctx context.Context, userName, password, inviteCode string) (uint64, error) {
	// 1. 校验邀请码并获取邀请人ID（无效或已过期直接拒绝）
	inviteKey := fmt.Sprintf("invite_code:%s", inviteCode)
	inviterIdValue, err := g.Redis().Get(ctx, inviteKey)
	if err != nil {
		return 0, gerror.Wrap(err, "查询邀请码失败")
	}
	if inviterIdValue.IsNil() {
		return 0, gerror.NewCode(gcode.CodeValidationFailed, "邀请码无效或已过期")
	}
	inviterId := inviterIdValue.Uint64()
	if inviterId == 0 {
		return 0, gerror.NewCode(gcode.CodeValidationFailed, "邀请码无效或已过期")
	}

	// 2. 用户名唯一性检查（软删除自动过滤，已删除账号的用户名可复用）
	count, err := dao.DcUser.Ctx(ctx).Where(do.DcUser{UserName: userName}).Count()
	if err != nil {
		return 0, gerror.Wrap(err, "查询用户名失败")
	}
	if count > 0 {
		return 0, gerror.NewCode(gcode.CodeValidationFailed, "用户名已被注册")
	}

	// 3. 密码 bcrypt 哈希
	hashedPassword, err := coding.HashPassword(password)
	if err != nil {
		return 0, gerror.Wrap(err, "密码加密失败")
	}

	// 4. 写入用户表（注册暂不采集邮箱，占位空串；等级与状态使用默认值 1）
	id, err := dao.DcUser.Ctx(ctx).Data(do.DcUser{
		UserName:  userName,
		Email:     "",
		Password:  hashedPassword,
		InviterId: inviterId,
		UserRank:  1,
		Status:    1,
	}).InsertAndGetId()
	if err != nil {
		return 0, gerror.Wrap(err, "创建用户失败")
	}

	// 5. 邀请码一次性使用，注册成功后立即删除；删除失败不影响注册结果
	if _, err = g.Redis().Del(ctx, inviteKey); err != nil {
		g.Log().Warningf(ctx, "注册成功但删除邀请码失败, key=%s: %v", inviteKey, err)
	}

	return uint64(id), nil
}
