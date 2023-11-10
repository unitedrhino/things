package firmwaremanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareDeviceInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
}

func NewOtaFirmwareDeviceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareDeviceInfoLogic {
	return &OtaFirmwareDeviceInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 获取固件包对应设备版本列表
func (l *OtaFirmwareDeviceInfoLogic) OtaFirmwareDeviceInfo(in *dm.OtaFirmwareDeviceInfoReq) (*dm.OtaFirmwareDeviceInfoResp, error) {
	var fDB = relationDB.NewOtaFirmwareRepo(l.ctx)
	di, err := fDB.FindOne(l.ctx, in.FirmwareID)
	if err != nil {
		return nil, err
	}
	result, err := l.DiDB.CountGroupByField(l.ctx, relationDB.DeviceFilter{ProductID: di.ProductID}, "version")
	if err != nil {
		return nil, err
	}
	versionString := ""
	for version, _ := range result {
		versionString = version + ","
	}
	//l.PiDB.FindByFilter()
	//TOD 查询product
	return &dm.OtaFirmwareDeviceInfoResp{Versions: versionString}, nil
}
