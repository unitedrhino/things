package sipmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/vidsip/internal/logic/common"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsip/internal/svc"
	"github.com/i-Things/things/src/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipDeviceUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipDeviceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipDeviceUpdateLogic {
	return &SipDeviceUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新GB28181设备
func (l *SipDeviceUpdateLogic) SipDeviceUpdate(in *sip.SipDevUpdateReq) (*sip.Response, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewSipDevicesRepo(l.ctx)

	po, err := deviceRepo.FindOneByFilter(l.ctx, db.SipDevicesFilter{
		DeviceIDs: []string{in.DeviceID},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find ID:" + string(in.DeviceID))
		}
		return nil, err
	}
	common.UpdatSipDeviceDB(po, in)
	if err := deviceRepo.Update(l.ctx, po); err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s req=%v err=%v", utils.FuncName(), po.DeviceID, er)
		return nil, er
	}
	return &sip.Response{}, nil
}
