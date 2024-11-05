package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	commonschemalogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/schemamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaUpdateLogic {
	return &DeviceSchemaUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新设备物模型
func (l *DeviceSchemaUpdateLogic) DeviceSchemaUpdate(in *dm.DeviceSchema) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithAllProject(l.ctx)

	po, err := relationDB.NewDeviceSchemaRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceSchemaFilter{
		ProductID: in.ProductID, DeviceName: in.DeviceName, Identifiers: []string{in.Identifier},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("标识符不存在:" + in.Identifier)
		}
		return nil, err
	}
	newPo := utils.Copy[relationDB.DmDeviceSchema](in)
	newPo.ID = po.ID
	if in.Affordance != nil && po.Tag == schema.TagCustom {
		po.Affordance = newPo.Affordance
	}
	if in.Name != nil {
		po.Name = newPo.Name
	}
	if in.Desc != nil {
		po.Desc = newPo.Desc
	}
	if in.Required != 0 {
		po.Required = newPo.Required
	}
	if in.IsCanSceneLinkage != 0 && po.Tag == schema.TagCustom {
		po.IsCanSceneLinkage = newPo.IsCanSceneLinkage
	}
	if in.FuncGroup != 0 && po.Tag == schema.TagCustom {
		po.FuncGroup = newPo.FuncGroup
	}
	if in.ControlMode != 0 && po.Tag == schema.TagCustom {
		po.ControlMode = newPo.ControlMode
	}
	if in.UserPerm != 0 && po.Tag == schema.TagCustom {
		po.UserPerm = newPo.UserPerm
	}

	if in.IsHistory != 0 && po.Tag == schema.TagCustom {
		po.IsHistory = newPo.IsHistory
	}

	if in.Order != 0 {
		po.Order = newPo.Order
	}

	if in.ExtendConfig != "" {
		po.ExtendConfig = newPo.ExtendConfig
	}
	if err := commonschemalogic.CheckAffordance(&newPo.DmSchemaCore); err != nil {
		return nil, err
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	//清除缓存
	err = l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	if err != nil {
		return nil, err
	}
	return &dm.Empty{}, nil
}
