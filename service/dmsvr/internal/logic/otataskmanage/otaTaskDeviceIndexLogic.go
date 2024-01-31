package otataskmanagelogic

import (
	"context"

	"gitee.com/i-Things/core/shared/def"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceIndexLogic {
	return &OtaTaskDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 升级批次详情列表
func (l *OtaTaskDeviceIndexLogic) OtaTaskDeviceIndex(in *dm.OtaTaskDeviceIndexReq) (*dm.OtaTaskDeviceIndexResp, error) {
	var (
		info     []*dm.OtaTaskDeviceInfo
		size     int64
		page     int64 = 1
		pageSize int64 = 20
		err      error
		otDB     = relationDB.NewOtaTaskDevicesRepo(l.ctx)
	)
	if !(in.Page == nil || in.Page.Page == 0 || in.Page.Size == 0) {
		page = in.Page.Page
		pageSize = in.Page.Size
	}

	size, err = otDB.CountByFilter(
		l.ctx, relationDB.OtaTaskDevicesFilter{
			FirmwareID: in.FirmwareID,
			TaskUid:    in.TaskUid,
			DeviceName: in.DeviceName,
			Status:     int64(in.Status),
		})
	if err != nil {
		return nil, err
	}
	di, err := otDB.FindByFilter(
		l.ctx, relationDB.OtaTaskDevicesFilter{
			FirmwareID: in.FirmwareID,
			TaskUid:    in.TaskUid,
			DeviceName: in.DeviceName,
			Status:     int64(in.Status),
		}, &def.PageInfo{Size: pageSize, Page: page})
	if err != nil {
		return nil, err
	}
	info = make([]*dm.OtaTaskDeviceInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToOtaTaskDeviceInfo(v))
	}
	return &dm.OtaTaskDeviceIndexResp{List: info, Total: size}, nil
}
