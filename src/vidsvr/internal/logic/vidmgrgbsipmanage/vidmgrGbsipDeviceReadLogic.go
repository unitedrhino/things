package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipDeviceReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipDeviceReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipDeviceReadLogic {
	return &VidmgrGbsipDeviceReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取GB28181设备详情
func (l *VidmgrGbsipDeviceReadLogic) VidmgrGbsipDeviceRead(in *vid.VidmgrGbsipDeviceReadReq) (*vid.VidmgrGbsipDevice, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewVidmgrDevicesRepo(l.ctx)

	filter := db.VidmgrDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	}
	device, err := deviceRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), device.DeviceID, er)
		return nil, er
	}
	return common.ToVidmgrGbsipDeviceRpc(device), nil
}
