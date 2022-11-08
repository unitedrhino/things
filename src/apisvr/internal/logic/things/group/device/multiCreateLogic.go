package device

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.GroupDeviceCreateReq) error {

	m := make([]*dm.DeviceInfoReadReq, 0, len(req.List))
	for _, v := range req.List {
		m = append(m, &dm.DeviceInfoReadReq{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	_, err := l.svcCtx.DeviceG.GroupDeviceCreate(l.ctx, &dm.GroupDeviceCreateReq{GroupID: req.GroupID, List: m})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.DeviceGroup Create req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}

	return nil
}
