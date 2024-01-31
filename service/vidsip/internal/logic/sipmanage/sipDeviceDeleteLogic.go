package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipDeviceDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipDeviceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipDeviceDeleteLogic {
	return &SipDeviceDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除GB28181设备
func (l *SipDeviceDeleteLogic) SipDeviceDelete(in *sip.SipDevDeleteReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	deviceRepo := db.NewSipDevicesRepo(l.ctx)
	channelRepo := db.NewSipChannelsRepo(l.ctx)
	filter := db.SipDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	}
	do, err := deviceRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.DeviceID))
		}
		return nil, err
	}

	filterChn := db.SipChannelsFilter{
		DeviceIDs: []string{do.DeviceID},
	}
	errChn := channelRepo.DeleteByFilter(l.ctx, filterChn)
	if errChn != nil {
		l.Errorf("%s.DeleteChn err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	//删除设备同时也要删除通道信息
	err = deviceRepo.DeleteByFilter(l.ctx, filter)
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return &sip.Response{}, nil
}
