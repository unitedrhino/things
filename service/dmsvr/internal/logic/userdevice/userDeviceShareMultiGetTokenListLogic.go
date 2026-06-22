package userdevicelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareMultiGetTokenListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareMultiGetTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareMultiGetTokenListLogic {
	return &UserDeviceShareMultiGetTokenListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserDeviceShareMultiGetTokenList 获取当前登录用户生成的批量分享 Token 列表
// 自动过滤已过期的 Token（惰性清理）
func (l *UserDeviceShareMultiGetTokenListLogic) UserDeviceShareMultiGetTokenList(in *dm.Empty) (*dm.UserDeviceShareMultiGetTokenListResp, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if uc == nil {
		return nil, errors.NotLogin
	}

	tenantCode := uc.TenantCode
	list, err := l.svcCtx.UserMultiDeviceShare.GetList(l.ctx, tenantCode, uc.UserID)
	if err != nil {
		return nil, err
	}

	var result []*dm.UserDeviceShareMultiListItem
	for _, item := range list {
		result = append(result, &dm.UserDeviceShareMultiListItem{
			ShareToken:  item.Token,
			DeviceCount: int64(len(item.Info.Devices)),
			CreatedTime: item.Info.CreatedTime,
			ExpTime:     item.Info.ExpTime,
			AuthType:    item.Info.AuthType,
			UseBy:       item.Info.UseBy,
			Desc:        item.Info.Desc,
		})
	}

	return &dm.UserDeviceShareMultiGetTokenListResp{
		List:  result,
		Total: int64(len(result)),
	}, nil
}
