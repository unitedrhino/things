package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolServiceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolServiceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolServiceIndexLogic {
	return &ProtocolServiceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProtocolServiceIndexLogic) ProtocolServiceIndex(in *dm.ProtocolServiceIndexReq) (*dm.ProtocolServiceIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	var (
		size int64
		err  error
		piDB = relationDB.NewProtocolServiceRepo(l.ctx)
	)

	filter := relationDB.ProtocolServiceFilter{
		Code: in.Code,
	}
	size, err = piDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	di, err := piDB.FindByFilter(l.ctx, filter,
		logic.ToPageInfo(in.Page),
	)
	if err != nil {
		return nil, err
	}

	return &dm.ProtocolServiceIndexResp{List: utils.CopySlice[dm.ProtocolService](di), Total: size}, nil
}
