package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/vidsip/internal/media"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipDeviceCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipDeviceCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipDeviceCreateLogic {
	return &SipDeviceCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新建GB28181设备
func (l *SipDeviceCreateLogic) SipDeviceCreate(in *sip.SipDevCreateReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewSipDevicesRepo(l.ctx)
	device := &db.SipDevices{
		DeviceID: in.DeviceID,
		VidmgrID: in.VidmgrID,
		//DeviceID: fmt.Sprintf("%s%06d", media.SipInfo.DID, media.SipInfo.DNUM+1),
		Region:    media.SipInfo.Region,
		Name:      in.Name,
		PWD:       in.PWD,
		MediaIP:   utils.InetAtoN(in.MediaIP),
		MediaPort: in.MediaPort,
	}
	if device.Name == "" {
		device.Name = device.DeviceID
	}

	if err := deviceRepo.Insert(l.ctx, device); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), device.DeviceID, er)
		return nil, errors.MediaSipDevCreateError.AddDetail(er)
	}

	return &sip.Response{}, nil
}
