package msg

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbnormalLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取设备异常日志
func NewAbnormalLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbnormalLogIndexLogic {
	return &AbnormalLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AbnormalLogIndexLogic) AbnormalLogIndex(req *types.DeviceMsgAbnormalLogIndexReq) (resp *types.DeviceMsgAbnormalLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.AbnormalLogIndex(l.ctx, utils.Copy[dm.AbnormalLogIndexReq](req))
	if err != nil {
		return nil, err
	}

	return &types.DeviceMsgAbnormalLogIndexResp{
		PageResp: logic.ToPageResp(req.Page, dmResp.Total),
		List:     utils.CopySlice[types.DeviceMsgAbnormalLogInfo](dmResp.List),
	}, nil
}
