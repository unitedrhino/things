package onlineCheck

import (
	"context"
	"github.com/i-Things/things/service/dgsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckEvent struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	ctx context.Context
}

func NewOnlineCheckEvent(svcCtx *svc.ServiceContext, ctx context.Context) *CheckEvent {
	return &CheckEvent{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (o *CheckEvent) Check() error {
	//o.svcCtx.DeviceM.DeviceInfoIndex(o.ctx, &dm.DeviceInfoIndexReq{
	//	Page: &dm.PageInfo{
	//		Page:   1,
	//		Size:   1000,
	//		Orders: []*dm.PageInfo_OrderBy{{Filed: "isOnline", Sort: 1}},
	//	},
	//	IsOnline:          0,
	//	ProductCategoryID: 0,
	//	Devices:           nil,
	//	IsShared:          0,
	//	TenantCode:        "",
	//	Versions:          nil,
	//})
	//isOnline, err := o.svcCtx.PubDev.CheckIsOnline(o.ctx, info.ClientID)
	//if err != nil {
	//	return err
	//}
	return nil
}
