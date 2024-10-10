package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/dmExport"
	"gitee.com/unitedrhino/things/service/udsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoIndexLogic {
	return &SceneInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoIndexLogic) SceneInfoIndex(in *ud.SceneInfoIndexReq) (*ud.SceneInfoIndexResp, error) {
	if in.Tag == "deviceTiming" { //单设备定时
		uc := ctxs.GetUserCtx(l.ctx)
		err := dmExport.AccessPerm(l.ctx, l.svcCtx.DeviceCache, l.svcCtx.UserShareCache, def.AuthRead, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		}, "deviceTiming")
		if err != nil {
			return nil, err
		}
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		if uc.ProjectID != di.ProjectID {
			uc.ProjectID = di.ProjectID
			uc.IsAdmin = true
		}
	}
	f := relationDB.SceneInfoFilter{AreaID: in.AreaID, IsCommon: in.IsCommon, Tag: in.Tag, Status: in.Status, Name: in.Name, DeviceMode: in.DeviceMode,
		Type: in.Type, HasActionType: in.HasActionType, IDs: in.SceneIDs, ProductID: in.ProductID, DeviceName: in.DeviceName}
	list, err := relationDB.NewSceneInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "createdTime",
		Sort:  stores.OrderDesc,
	}))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewSceneInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}

	return &ud.SceneInfoIndexResp{List: PoToSceneInfoPbs(l.ctx, l.svcCtx, list), Total: total}, nil
}
