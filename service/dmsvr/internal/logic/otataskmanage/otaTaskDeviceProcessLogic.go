package otataskmanagelogic

import (
	"context"

	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceProcessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceProcessLogic {
	return &OtaTaskDeviceProcessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 升级进度上报
func (l *OtaTaskDeviceProcessLogic) OtaTaskDeviceProcess(in *dm.OtaTaskDeviceProcessReq) (*dm.OtaCommonResp, error) {
	var otDB = relationDB.NewOtaTaskDevicesRepo(l.ctx)
	di, err := otDB.FindOne(l.ctx, in.ID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind.AddDetailf("not find ota task info|id=%s",
				in.ID)
		}
		return nil, errors.Database.AddDetail(err)
	}
	l.ChangeOtaTaskDevice(di, in)
	err = otDB.Update(l.ctx, di)
	if err != nil {
		l.Errorf("OtaTaskDeviceProcess.OtaTaskDeviceInfo.Update err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.OtaCommonResp{}, nil
}
func (l *OtaTaskDeviceProcessLogic) ChangeOtaTaskDevice(old *relationDB.DmOtaTaskDevices, data *dm.OtaTaskDeviceProcessReq) {
	if data.Desc != "" {
		old.Desc = data.Desc
	}
	if data.Step != 0 {
		old.Step = data.Step
		if data.Step < 0 {
			old.Status = 601 //升级失败
		} else if data.Step > 0 {
			old.Status = 401 //升级中
		}
		//物联网平台收到了新的版本号上报后, 才会判定升级成功, 否则会认定升级失败
	}
	//TODO 待推送的状态怎么判断呢？ 重试的时候怎么重置
}
