package otataskmanagelogic

import (
	"context"

	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceCancleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceCancleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceCancleLogic {
	return &OtaTaskDeviceCancleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消单个设备的升级
func (l *OtaTaskDeviceCancleLogic) OtaTaskDeviceCancle(in *dm.OtaTaskDeviceCancleReq) (*dm.OtaCommonResp, error) {
	var otDB = relationDB.NewOtaTaskDevicesRepo(l.ctx)
	otd, err := otDB.FindOne(l.ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if otd.Status == 501 || otd.Status == 601 {
		return nil, errors.OtaCancleStatusError
	}
	otd.Status = 701
	err = otDB.Update(l.ctx, otd)
	return &dm.OtaCommonResp{}, err
}
