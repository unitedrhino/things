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

type VidmgrGbsipDeviceUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipDeviceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipDeviceUpdateLogic {
	return &VidmgrGbsipDeviceUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新GB28181设备
func (l *VidmgrGbsipDeviceUpdateLogic) VidmgrGbsipDeviceUpdate(in *vid.VidmgrGbsipDeviceUpdateReq) (*vid.Response, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewVidmgrDevicesRepo(l.ctx)

	po, err := deviceRepo.FindOneByFilter(l.ctx, db.VidmgrDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find ID:" + string(in.DeviceID))
		}
		return nil, err
	}
	common.UpdatVidmgrDeviceDB(po, in)
	if err := deviceRepo.Update(l.ctx, po); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), po.DeviceID, er)
		return nil, er
	}
	return &vid.Response{}, nil
}
