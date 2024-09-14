package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoDeleteLogic {
	return &SceneInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoDeleteLogic) SceneInfoDelete(in *ud.WithID) (*ud.Empty, error) {
	old, err := relationDB.NewSceneInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if old.Tag == "deviceTiming" { //单设备定时
		uc := ctxs.GetUserCtx(l.ctx)
		err := dmExport.AccessPerm(l.ctx, l.svcCtx.DeviceCache, l.svcCtx.UserShareCache, def.AuthReadWrite, devices.Core{
			ProductID:  old.ProductID,
			DeviceName: old.DeviceName,
		}, "deviceTiming")
		if err != nil {
			return nil, err
		}
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  old.ProductID,
			DeviceName: old.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		if uc.ProjectID != di.ProjectID {
			uc.ProjectID = di.ProjectID
			uc.IsAdmin = true
		}
	}
	err = relationDB.NewSceneInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &ud.Empty{}, err
}
