package devicemsglogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbnormalLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbnormalLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbnormalLogIndexLogic {
	return &AbnormalLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AbnormalLogIndexLogic) AbnormalLogIndex(in *dm.AbnormalLogIndexReq) (*dm.AbnormalLogIndexResp, error) {
	filter := deviceLog.AbnormalFilter{
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		Action:      in.Action,
		Type:        in.Type,
		AreaIDs:     in.AreaIDs,
		AreaID:      in.AreaID,
		AreaIDPath:  in.AreaIDPath,
		ProjectID:   in.ProjectID,
		BelongGroup: utils.CopyMap3[def.IDsInfo](in.BelongGroup),
	}
	if in.ProductID != "" && in.DeviceName != "" {
		_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		}, nil)
		if err != nil {
			return nil, err
		}
	} else {
		uc := ctxs.GetUserCtxNoNil(l.ctx)
		if uc.IsAdmin != true {
			return nil, errors.Parameter.AddMsg("请填写产品和设备")
		}
		if ctxs.IsRoot(l.ctx) != nil {
			filter.TenantCode = uc.TenantCode
		}
		if uc.ProjectID > def.NotClassified {
			filter.ProjectID = uc.ProjectID
		}
		if in.ProductCategoryID != 0 {
			pis, err := relationDB.NewProductInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.ProductFilter{CategoryIDs: []int64{in.ProductCategoryID}}, nil)
			if err != nil {
				return &dm.AbnormalLogIndexResp{}, err
			}
			for _, pi := range pis {
				filter.ProductIDs = append(filter.ProductIDs, pi.ProductID)
			}
		}
	}

	page := def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if !uc.IsAdmin {
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		if di.LastBind*1000 > page.TimeStart {
			page.TimeStart = di.LastBind * 1000
		}
	}
	logs, err := l.svcCtx.AbnormalRepo.GetDeviceLog(l.ctx, filter, page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	total, err := l.svcCtx.AbnormalRepo.GetCountLog(l.ctx, filter, page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*dm.AbnormalLogInfo
	for _, v := range logs {
		data = append(data, ToDataAbnormalLogIndex(v))
	}
	return &dm.AbnormalLogIndexResp{List: data, Total: total}, nil
}
