package otaupgradetaskmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceProcessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaUpgradeTaskRepo
}

func NewOtaTaskDeviceProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceProcessLogic {
	return &OtaTaskDeviceProcessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
	}
}

// 根据设备的name和productId，查询task
func (l *OtaTaskDeviceProcessLogic) OtaTaskDeviceProcess(in *dm.OtaTaskDeviceProcessReq) (*dm.Response, error) {
	filter := relationDB.OtaUpgradeTaskFilter{
		ProductId:  in.ProductId,
		DeviceName: in.DeviceName,
		Module:     in.Module,
	}
	task, err := l.OtDB.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind.AddDetailf("not find ota task info|id=%s",
				in.ID)
		}
		return nil, errors.Database.AddDetail(err)
	}
	//更新一下进度
	//todo
	//状态啥的需要改变一下
	task.Step = in.Step
	err = l.OtDB.Update(l.ctx, task)
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfo OtaTaskDeviceProcess failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
