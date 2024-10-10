package device

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiCreateLogic {
	return &MultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.GroupDeviceMultiSaveReq) error {
	m := make([]*dm.DeviceCore, 0, len(req.List))
	for _, v := range req.List {
		m = append(m, &dm.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	_, err := l.svcCtx.DeviceG.GroupDeviceMultiCreate(l.ctx, &dm.GroupDeviceMultiSaveReq{GroupID: req.GroupID, List: m})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GroupDeviceMultiCreate req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
