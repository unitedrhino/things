package otamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func ToFirmwareInfoPb(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.DmOtaFirmwareInfo) *dm.OtaFirmwareInfo {
	result := utils.Copy[dm.OtaFirmwareInfo](in)
	pi, _ := svcCtx.ProductCache.GetData(ctx, in.ProductID)
	if pi != nil {
		result.ProductName = pi.ProductName
	}
	utils.CopyE(&result.FileList, &in.Files)
	return result
}

func ToFirmwareDeviceInfo(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.DmOtaFirmwareDevice) *dm.OtaFirmwareDeviceInfo {
	var result = dm.OtaFirmwareDeviceInfo{CreatedTime: in.CreatedTime.Unix(), UpdatedTime: in.UpdatedTime.Unix()}
	utils.CopyE(&result, &in)
	pi, _ := svcCtx.ProductCache.GetData(ctx, in.ProductID)
	if pi != nil {
		result.ProductName = pi.ProductName
	}
	return &result
}
