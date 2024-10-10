package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

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
	po, err := SceneInfoRead(l.ctx, l.svcCtx, in.Id, def.AuthRead)
	return PoToSceneInfoPb(l.ctx, l.svcCtx, po), err
}

func SceneInfoRead(ctx context.Context, svcCtx *svc.ServiceContext, id int64, perm def.AuthType) (*relationDB.UdSceneInfo, error) {
	po, err := relationDB.NewSceneInfoRepo(ctx).FindOne(ctx, id)
	//需要支持分享的设备
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}
	if po == nil { //可能是共享的设备定时
		ctx = ctxs.WithAllProject(ctx)
		po, err = relationDB.NewSceneInfoRepo(ctx).FindOne(ctx, id)
		if err != nil {
			return nil, err
		}
		if po.Tag != "deviceTiming" { //单设备定时
			return nil, errors.NotFind
		}
		err := dmExport.AccessPerm(ctx, svcCtx.DeviceCache, svcCtx.UserShareCache, perm, devices.Core{
			ProductID:  po.ProductID,
			DeviceName: po.DeviceName,
		}, "deviceTiming")
		if err != nil {
			return nil, err
		}
	}
	return po, nil
}
