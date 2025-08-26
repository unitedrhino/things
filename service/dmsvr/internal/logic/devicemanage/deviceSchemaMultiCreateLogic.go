package devicemanagelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaMultiCreateLogic {
	return &DeviceSchemaMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func ruleCheck(ctx context.Context, svcCtx *svc.ServiceContext, in *dm.DeviceSchemaMultiCreateReq) (ret []*relationDB.DmDeviceSchema, err error) {
	_, err = relationDB.NewDeviceInfoRepo(ctx).FindOneByFilter(ctx,
		relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("找不到该设备:%v", in.DeviceName)
		}
		return nil, err
	}
	s, err := svcCtx.DeviceSchemaRepo.GetData(ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		return nil, err
	}
	for _, v := range in.List {
		var checkOptions []logic.CheckOption
		po := utils.Copy[relationDB.DmDeviceSchema](v)
		po.ProductID = in.ProductID
		po.DeviceName = in.DeviceName
		switch schema.AffordanceType(v.Type) {
		case schema.AffordanceTypeProperty:
			_, ok := s.Property[v.Identifier]
			if ok {
				continue
			}
			if po.Tag == schema.TagOptional { //通用物模型里选取的
				po.Tag = schema.TagDeviceOptional
			}
			if po.Tag != schema.TagDeviceOptional {
				checkOptions = append(checkOptions, func(do any) error {
					s := do.(*schema.Property)
					if utils.SliceIn(s.Define.Type, schema.DataTypeArray, schema.DataTypeStruct) {
						return errors.Parameter.AddMsgf("自定义物模型中不支持数组或结构体:%v", v.Identifier)
					}
					return nil
				})
				po.Tag = schema.TagDeviceCustom
			}

		case schema.AffordanceTypeAction:
			if v.Tag != schema.TagDeviceOptional {
				continue
			}
			_, ok := s.Action[v.Identifier]
			if ok {
				continue
			}
			po.Tag = schema.TagDeviceOptional
		case schema.AffordanceTypeEvent:
			if v.Tag != schema.TagDeviceOptional {
				continue
			}
			_, ok := s.Event[v.Identifier]
			if ok {
				continue
			}
			po.Tag = schema.TagDeviceOptional
		}
		if po.Tag == schema.TagDeviceOptional {
			err := func() error {
				//如果导入的是通用物模型
				var cs *relationDB.DmCommonSchema
				cs, err = relationDB.NewCommonSchemaRepo(ctx).FindOneByFilter(ctx, relationDB.CommonSchemaFilter{Identifiers: []string{v.Identifier}})
				if err != nil {
					return err
				}
				if cs == nil { //如果通用物模型里面没有,则变成自定义物模型
					po.Tag = schema.TagDeviceCustom
					checkOptions = append(checkOptions, func(do any) error {
						s := do.(*schema.Property)
						if utils.SliceIn(s.Define.Type, schema.DataTypeArray, schema.DataTypeStruct) {
							return errors.Parameter.AddMsgf("自定义物模型中不支持数组或结构体:%v", v.Identifier)
						}
						return nil
					})
					return nil
				}
				po.IsCanSceneLinkage = cs.IsCanSceneLinkage
				po.FuncGroup = cs.FuncGroup
				po.ControlMode = cs.ControlMode
				po.UserPerm = cs.UserPerm
				po.RecordMode = cs.RecordMode
				po.IsPassword = cs.IsPassword
				if po.Name == "" {
					po.Name = cs.Name
				}
				if po.Required == 0 {
					po.Required = cs.Required
				}
				if po.IsCanSceneLinkage == 0 {
					po.IsCanSceneLinkage = cs.IsCanSceneLinkage
				}
				if po.FuncGroup == 0 {
					po.FuncGroup = cs.FuncGroup
				}
				if po.ControlMode == 0 {
					po.ControlMode = cs.ControlMode
				}
				if po.UserPerm != 0 {
					po.UserPerm = cs.UserPerm
				}
				if po.RecordMode == 0 {
					po.RecordMode = cs.RecordMode
				}
				if po.Order == 0 {
					po.Order = cs.Order
				}
				if po.IsPassword == 0 {
					po.IsPassword = cs.IsPassword
				}
				if po.ExtendConfig == "" {
					po.ExtendConfig = cs.ExtendConfig
				}
				return nil
			}()
			if err != nil {
				return nil, err
			}
		}

		if err = logic.CheckAffordance(po.Identifier, &po.DmSchemaCore, nil, checkOptions...); err != nil {
			return nil, err
		}
		ret = append(ret, po)
	}
	return
}

// 批量新增物模型,只新增没有的,已有的不处理
func (l *DeviceSchemaMultiCreateLogic) DeviceSchemaMultiCreate(in *dm.DeviceSchemaMultiCreateReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithAllProject(l.ctx)
	pos, err := ruleCheck(l.ctx, l.svcCtx, in)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).MultiInsert(l.ctx, pos)
	if err != nil {
		l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	}
	return &dm.Empty{}, err
}
