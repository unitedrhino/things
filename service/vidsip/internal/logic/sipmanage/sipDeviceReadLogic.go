package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsip/internal/logic/common"
	"github.com/i-Things/things/service/vidsip/internal/repo/relationDB"

	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipDeviceReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipDeviceReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipDeviceReadLogic {
	return &SipDeviceReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取GB28181设备详情
func (l *SipDeviceReadLogic) SipDeviceRead(in *sip.SipDevReadReq) (*sip.SipDevice, error) {
	// todo: add your logic here and delete this line
	deviceRepo := relationDB.NewSipDevicesRepo(l.ctx)

	filter := relationDB.SipDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	}
	device, err := deviceRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), device.DeviceID, er)
		return nil, er
	}
	return common.ToSipDeviceRpc(device), nil
}
