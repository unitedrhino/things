package device

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateLogic {
	return &MultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUpdateLogic) MultiUpdate(req *types.GroupDeviceMultiSaveReq) error {
	m := make([]*dm.DeviceCore, 0, len(req.List))
	for _, v := range req.List {
		m = append(m, &dm.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	_, err := l.svcCtx.DeviceG.GroupDeviceMultiUpdate(l.ctx, &dm.GroupDeviceMultiSaveReq{GroupID: req.GroupID, List: m})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GroupDeviceMultiCreate req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
