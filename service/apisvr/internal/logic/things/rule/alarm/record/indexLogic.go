package record

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.AlarmRecordIndexReq) (resp *types.AlarmRecordIndexResp, err error) {
	ret, err := l.svcCtx.Rule.AlarmRecordIndex(l.ctx, &ud.AlarmRecordIndexReq{
		AlarmID:    req.AlarmID,
		Page:       logic.ToUdPageRpc(req.Page),
		TimeRange:  logic.ToUdTimeRangeRpc(req.TimeRange),
		DealStatus: req.DealStatus,
		AlarmName:  req.AlarmName,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmRecordIndex req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return utils.Copy[types.AlarmRecordIndexResp](ret), nil
}
