package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoReadLogic {
	return &SceneInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoReadLogic) SceneInfoRead(in *ud.WithID) (*ud.SceneInfo, error) {
	po, err := relationDB.NewSceneInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	//需要支持分享的设备
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}
	if po == nil { //可能是共享的设备定时
		l.ctx = ctxs.WithAllProject(l.ctx)
		po, err = relationDB.NewSceneInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
		if po.Tag != "deviceTiming" { //单设备定时
			return nil, errors.NotFind
		}
		err := dmExport.AccessPerm(l.ctx, l.svcCtx.DeviceCache, l.svcCtx.UserShareCache, def.AuthRead, devices.Core{
			ProductID:  po.ProductID,
			DeviceName: po.DeviceName,
		}, "deviceTiming")
		if err != nil {
			return nil, err
		}

	}
	return PoToSceneInfoPb(l.ctx, l.svcCtx, po), err
}
