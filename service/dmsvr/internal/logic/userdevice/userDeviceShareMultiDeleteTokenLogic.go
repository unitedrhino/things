package userdevicelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareMultiDeleteTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareMultiDeleteTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareMultiDeleteTokenLogic {
	return &UserDeviceShareMultiDeleteTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserDeviceShareMultiDeleteToken 删除批量分享 Token
// 仅 Token 创建者本人可删除
func (l *UserDeviceShareMultiDeleteTokenLogic) UserDeviceShareMultiDeleteToken(in *dm.UserDeviceShareMultiDeleteTokenReq) (*dm.Empty, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if uc == nil {
		return nil, errors.NotLogin
	}

	tenantCode := uc.TenantCode
	info, err := l.svcCtx.UserMultiDeviceShare.GetData(l.ctx, in.ShareToken)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("分享 Token 不存在或已过期")
		}
		return nil, err
	}

	// 验证所有权：只能删除自己生成的分享 Token
	if info.UserID != uc.UserID {
		return nil, errors.Permissions.AddMsg("只能删除自己生成的分享 Token")
	}

	err = l.svcCtx.UserMultiDeviceShare.DeleteToken(l.ctx, tenantCode, uc.UserID, in.ShareToken)
	if err != nil {
		return nil, err
	}

	return &dm.Empty{}, nil
}
