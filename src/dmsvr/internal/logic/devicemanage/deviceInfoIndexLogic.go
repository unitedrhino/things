package devicemanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/spf13/cast"
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
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	var (
		info []*dm.DeviceInfo
		size int64
		err  error
	)

	position := "POINT(0 0)"
	if in.Position != nil {
		position = fmt.Sprintf("POINT(%s)",
			cast.ToString(in.Position.Longitude)+" "+cast.ToString(in.Position.Latitude))
	}
	filter := mysql.DeviceFilter{
		ProductID:   in.ProductID,
		AreaIDs:     in.AreaIDs,
		DeviceName:  in.DeviceName,
		Tags:        in.Tags,
		Range:       in.Range,
		Position:    position,
		DeviceAlias: in.DeviceAlias,
	}

	size, err = l.svcCtx.DeviceInfo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := l.svcCtx.DeviceInfo.FindByFilter(l.ctx, filter,
		logic.ToPageInfoWithDefault(in.Page, &def.PageInfo{
			Page: 1, Size: 20,
			Orders: []def.OrderBy{{"createdTime", def.OrderDesc}, {"productID", def.OrderDesc}},
		}),
	)
	if err != nil {
		return nil, err
	}

	info = make([]*dm.DeviceInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToDeviceInfo(v))
	}

	return &dm.DeviceInfoIndexResp{List: info, Total: size}, nil
}
