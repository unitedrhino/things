package share

import (
	"context"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiGetTokenListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取批量分享 Token 列表
func NewMultiGetTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiGetTokenListLogic {
	return &MultiGetTokenListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiGetTokenListLogic) MultiGetTokenList() (resp *types.UserDeviceShareMultiGetTokenListResp, err error) {
	ret, err := l.svcCtx.UserDevice.UserDeviceShareMultiGetTokenList(l.ctx, &dm.Empty{})
	if err != nil {
		return nil, err
	}
	var list []*types.UserDeviceShareMultiListItem
	for _, v := range ret.List {
		list = append(list, &types.UserDeviceShareMultiListItem{
			ShareToken:  v.ShareToken,
			DeviceCount: v.DeviceCount,
			CreatedTime: v.CreatedTime,
			ExpTime:     v.ExpTime,
			AuthType:    v.AuthType,
			UseBy:       v.UseBy,
		})
	}
	return &types.UserDeviceShareMultiGetTokenListResp{
		List:  list,
		Total: ret.Total,
	}, nil
}
