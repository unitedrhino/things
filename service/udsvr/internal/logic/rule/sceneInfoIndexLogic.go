package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/udsvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

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
		}
	}
	f := relationDB.SceneInfoFilter{AreaID: in.AreaID, IsCommon: in.IsCommon, Tag: in.Tag, Status: in.Status, Name: in.Name,
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
