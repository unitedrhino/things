package otamanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
)

func ToFirmwareInfoPb(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.DmOtaFirmwareInfo) *dm.OtaFirmwareInfo {
	var result = dm.OtaFirmwareInfo{FileList: []*dm.FirmwareFile{}, CreatedTime: in.CreatedTime.Unix()}
	copier.Copy(&result, &in)
	pi, _ := svcCtx.ProductCache.GetData(ctx, in.ProductID)
	if pi != nil {
		result.ProductName = pi.ProductName
	}
	copier.Copy(&result.FileList, &in.Files)
	return &result
}

func ToFirmwareDeviceInfo(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.DmOtaFirmwareDevice) *dm.OtaFirmwareDeviceInfo {
	var result = dm.OtaFirmwareDeviceInfo{CreatedTime: in.CreatedTime.Unix(), UpdatedTime: in.UpdatedTime.Unix()}
	copier.Copy(&result, &in)
	pi, _ := svcCtx.ProductCache.GetData(ctx, in.ProductID)
	if pi != nil {
		result.ProductName = pi.ProductName
	}
	return &result
}
