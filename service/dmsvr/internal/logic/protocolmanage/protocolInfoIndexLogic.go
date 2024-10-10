package protocolmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProtocolInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolInfoIndexLogic {
	return &ProtocolInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 协议列表
func (l *ProtocolInfoIndexLogic) ProtocolInfoIndex(in *dm.ProtocolInfoIndexReq) (*dm.ProtocolInfoIndexResp, error) {
	var (
		info []*dm.ProtocolInfo
		size int64
		err  error
		piDB = relationDB.NewProtocolInfoRepo(l.ctx)
	)

	filter := relationDB.ProtocolInfoFilter{
		Name:          in.Name,
		Code:          in.Code,
		TransProtocol: in.TransProtocol,
		NotCodes:      in.NotCodes,
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

	info = make([]*dm.ProtocolInfo, 0, len(di))
	for _, v := range di {
		info = append(info, logic.ToProtocolInfoPb(v))
	}
	return &dm.ProtocolInfoIndexResp{List: info, Total: size}, nil
}
