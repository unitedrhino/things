package info

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
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

func (l *ReadLogic) Read(req *types.WithIDOrCode) (resp *types.AppInfo, err error) {
	ret, err := l.svcCtx.AppRpc.AppInfoRead(l.ctx, system.ToSysReqWithIDCode(req))
	return ToAppInfoTypes(ret), err
}
