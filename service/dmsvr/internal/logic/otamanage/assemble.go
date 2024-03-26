package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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
