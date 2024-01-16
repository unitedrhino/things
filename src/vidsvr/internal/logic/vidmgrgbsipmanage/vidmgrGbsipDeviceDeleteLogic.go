package vidmgrgbsipmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipDeviceDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipDeviceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipDeviceDeleteLogic {
	return &VidmgrGbsipDeviceDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除GB28181设备
func (l *VidmgrGbsipDeviceDeleteLogic) VidmgrGbsipDeviceDelete(in *vid.VidmgrGbsipDeviceDeleteReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewVidmgrDevicesRepo(l.ctx)
	channelRepo := db.NewVidmgrChannelsRepo(l.ctx)
	filter := db.VidmgrDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	}
	do, err := deviceRepo.FindOneByFilter(l.ctx, filter)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("Channel not find ID:" + string(in.DeviceID))
		}
		return nil, err
	}

	filterChn := db.VidmgrChannelsFilter{
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
	return &vid.Response{}, nil
}
