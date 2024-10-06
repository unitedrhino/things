package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/i-Things/things/service/udsvr/internal/svc"
	"gitee.com/i-Things/things/service/udsvr/pb/ud"

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
	old, err := SceneInfoRead(l.ctx, l.svcCtx, in.Id, def.AuthReadWrite)
	if err != nil {
		return nil, err
	}
	if old.Tag == "deviceTiming" { //单设备定时
		uc := ctxs.GetUserCtx(l.ctx)
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
