package otaupgradetaskmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaFirmwareDeviceRepo
}

func NewOtaTaskConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskConfirmLogic {
	return &OtaTaskConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaFirmwareDeviceRepo(ctx),
	}
}

//
//// 批量确认，处于待确认状态的设备升级作业
//func (l *OtaTaskConfirmLogic) OtaTaskConfirm(in *dm.OTATaskConfirmReq) (*dm.Empty, error) {
//	filter := relationDB.OtaFirmwareDeviceFilter{
//		IDs: in.TaskIDs,
//	}
//	updateData := make(map[string]interface{})
//	err := l.OtDB.BatchUpdateField(l.ctx, filter, updateData)
//	if err != nil {
//		l.Errorf("%s.TaskInfo.TaskInfo Updates failure err=%+v", utils.FuncName(), err)
//		return nil, err
//	}
//	return &dm.Empty{}, nil
//}
