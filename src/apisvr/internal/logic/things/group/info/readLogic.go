package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.GroupInfoReadReq) (resp *types.GroupInfo, err error) {
	dgResp, err := l.svcCtx.DeviceG.GroupInfoRead(l.ctx, &dm.GroupInfoReadReq{GroupID: req.GroupID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.GroupInfo{
		GroupID:     dgResp.GroupID,
		GroupName:   dgResp.GroupName,
		CreatedTime: dgResp.CreatedTime,
		Desc:        dgResp.Desc,
		Tags:        toTagsType(dgResp.Tags),
	}, err
}
