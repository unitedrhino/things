package rulelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/ops"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRecordDealLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmRecordDealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRecordDealLogic {
	return &AlarmRecordDealLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type OpsParam struct {
	DeviceAlias string `json:"deviceAlias"`
	DeviceName  string `json:"deviceName"`
	ProductID   string `json:"productID"`
}

const (
	HandelIgnore    = 1 //忽略
	HandleSendOrder = 2 //派单
)

func (l *AlarmRecordDealLogic) AlarmRecordDeal(in *ud.AlarmRecordDealReq) (*ud.Empty, error) {
	po, err := relationDB.NewAlarmRecordRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if po.DealStatus != 1 {
		return &ud.Empty{}, errors.Parameter.AddMsg("只有告警中可以处理")
	}
	switch in.Handle {
	case HandelIgnore:
		po.DealStatus = scene.AlarmDealStatusIgnore
		err := relationDB.NewAlarmRecordRepo(l.ctx).Update(l.ctx, po)
		if err != nil {
			return nil, err
		}
	case HandleSendOrder:
		owo := sys.OpsWorkOrder{
			Type:      ops.WorkOrderTypeSceneAlarm,
			IssueDesc: "自动化告警创建",
			Params:    make(map[string]string),
		}
		if po.DeviceName != "" {
			di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{DeviceName: po.DeviceName, ProductID: po.ProductID})
			if err != nil {
				return nil, err
			}
			owo.AreaID = di.AreaID
			owo.Params["deviceAlias"] = di.DeviceAlias.GetValue()
			owo.Params["deviceName"] = di.DeviceName
			owo.Params["productID"] = di.ProductID
		}
		ret, err := l.svcCtx.Ops.OpsWorkOrderCreate(l.ctx, &owo)
		if err != nil {
			return nil, err
		}
		po.WorkOrderID = ret.Id
		po.DealStatus = scene.AlarmDealStatusInHand
		err = relationDB.NewAlarmRecordRepo(l.ctx).Update(l.ctx, po)
		if err != nil {
			return nil, err
		}
	}

	return &ud.Empty{}, nil
}
