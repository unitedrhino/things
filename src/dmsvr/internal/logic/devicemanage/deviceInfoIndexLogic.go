package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoIndexLogic {
	return &DeviceInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备信息列表
func (l *DeviceInfoIndexLogic) DeviceInfoIndex(in *dm.DeviceInfoIndexReq) (*dm.DeviceInfoIndexResp, error) {
	l.Infof("%s req=%+v",utils.FuncName(), in)
	var (
		info     []*dm.DeviceInfo
		size     int64
		page     int64 = 1
		pageSize int64 = 20
		err      error
	)
	if !(in.Page == nil || in.Page.Page == 0 || in.Page.Size == 0) {
		page = in.Page.Page
		pageSize = in.Page.Size
	}
	size, err = l.svcCtx.DmDB.GetDevicesCountByFilter(
		l.ctx, mysql.DeviceFilter{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
			Tags:       in.Tags,
		})
	if err != nil {
		return nil, err
	}
	di, err := l.svcCtx.DmDB.FindDevicesByFilter(
		l.ctx, mysql.DeviceFilter{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
			Tags:       in.Tags,
		}, def.PageInfo{Size: pageSize, Page: page})
	if err != nil {
		return nil, err
	}
	info = make([]*dm.DeviceInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToDeviceInfo(v))
	}
	return &dm.DeviceInfoIndexResp{List: info, Total: size}, nil
}
