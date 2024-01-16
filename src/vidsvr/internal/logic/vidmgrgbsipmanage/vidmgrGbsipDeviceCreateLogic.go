package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/internal/media"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipDeviceCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipDeviceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipDeviceCreateLogic {
	return &VidmgrGbsipDeviceCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新建GB28181设备
func (l *VidmgrGbsipDeviceCreateLogic) VidmgrGbsipDeviceCreate(in *vid.VidmgrGbsipDeviceCreateReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewVidmgrDevicesRepo(l.ctx)
	device := &db.VidmgrDevices{
		DeviceID: in.DeviceID,
		//DeviceID: fmt.Sprintf("%s%06d", media.SipInfo.DID, media.SipInfo.DNUM+1),
		Region: media.SipInfo.Region,
		Name:   in.Name,
		PWD:    in.PWD,
	}
	if device.Name == "" {
		device.Name = device.DeviceID
	}

	if err := deviceRepo.Insert(l.ctx, device); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), device.DeviceID, er)
		return nil, errors.MediaGbsipDevCreateError.AddDetail(er)
	}

	return &vid.Response{}, nil
}
