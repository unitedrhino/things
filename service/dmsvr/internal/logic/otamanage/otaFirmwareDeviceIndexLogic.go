package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaFirmwareDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareDeviceIndexLogic {
	return &OtaFirmwareDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询指定升级批次下的设备升级作业列表
func (l *OtaFirmwareDeviceIndexLogic) OtaFirmwareDeviceIndex(in *dm.OtaFirmwareDeviceIndexReq) (*dm.OtaFirmwareDeviceIndexResp, error) {
	//todo debug
	//if err := ctxs.IsRoot(l.ctx); err != nil {
	//	return nil, err
	//}
	l.ctx = ctxs.WithRoot(l.ctx)
	f := relationDB.OtaFirmwareDeviceFilter{FirmwareID: in.FirmwareID, JobID: in.JobID, DeviceName: in.DeviceName}
	repo := relationDB.NewOtaFirmwareDeviceRepo(l.ctx)
	total, err := repo.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := repo.FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		FirmwareID: in.FirmwareID,
		JobID:      in.JobID,
		DeviceName: in.DeviceName}, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	var (
		list = []*dm.OtaFirmwareDeviceInfo{}
	)
	for _, v := range pos {
		list = append(list, ToFirmwareDeviceInfo(l.ctx, l.svcCtx, v))
	}
	return &dm.OtaFirmwareDeviceIndexResp{List: list, Total: total}, nil
}
