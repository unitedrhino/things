package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/logic/devicemanage"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeviceIndexLogic {
	return &GroupDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分组设备信息列表
func (l *GroupDeviceIndexLogic) GroupDeviceIndex(in *dm.GroupDeviceIndexReq) (*dm.GroupDeviceIndexResp, error) {

	var list []*dm.DeviceInfo
	gd, total, err := l.svcCtx.GroupDB.IndexGD(l.ctx, &mysql.GroupDeviceFilter{
		Page:       logic.ToPageInfo(in.Page),
		GroupID:    in.GroupID,
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	for _, v := range gd {
		di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, v.ProductID, v.DeviceName)
		if err != nil {
			l.Errorf("%s.GroupInfo.DeviceInfoRead failure err=%+v", utils.FuncName(), err)
			continue
		}
		dd := devicemanagelogic.ToDeviceInfo(di)
		list = append(list, dd)
	}

	return &dm.GroupDeviceIndexResp{Total: total, List: list}, nil
}
