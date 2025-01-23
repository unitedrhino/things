package userdevicelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"github.com/hashicorp/go-uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareMultiCreateLogic {
	return &UserDeviceShareMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeviceShareMultiCreateLogic) UserDeviceShareMultiCreate(in *dm.UserDeviceShareMultiInfo) (*dm.UserDeviceShareMultiToken, error) {
	// 写入caches
	shareToken, _ := uuid.GenerateUUID()
	uc := ctxs.GetUserCtx(l.ctx)
	in.UserID = uc.UserID
	//判断是否有分享的权限
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(uc.ProjectID)})
	if err != nil {
		return nil, err
	}
	if pi.AdminUserID != uc.UserID && !uc.IsAdmin {
		pa := uc.ProjectAuth[pi.ProjectID]
		if pa.AuthType != def.AuthAdmin {
			if pa.Area == nil {
				return nil, errors.Permissions.AddMsg("只有管理员才能分享设备")
			}
			for _, d := range in.Devices {
				di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
					ProductID:  d.ProductID,
					DeviceName: d.DeviceName,
				})
				if err != nil {
					return nil, errors.Permissions.AddMsg("你分享了异常的设备")
				} else {
					if pa.Area[int64(di.AreaID)] != def.AuthAdmin {
						return nil, errors.Permissions.AddMsg("您无权分享所选的设备")
					}
				}
				d.DeviceAlias = di.DeviceAlias
				d.ProductName = di.ProductName
				d.ProductImg = di.ProductImg
			}
		}
	}
	for _, d := range in.Devices {
		//补全设备信息
		if d.ProductImg == "" {
			di, _ := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
				ProductID:  d.ProductID,
				DeviceName: d.DeviceName,
			})
			d.DeviceAlias = di.DeviceAlias
			d.ProductName = di.ProductName
			d.ProductImg = di.ProductImg
		}
	}
	l.svcCtx.UserMultiDeviceShare.SetData(l.ctx, shareToken, in)
	return &dm.UserDeviceShareMultiToken{ShareToken: shareToken}, nil
}
