package msg

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogIndexLogic {
	return &SendLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SendLogIndexLogic) SendLogIndex(req *types.DeviceMsgSendLogIndexReq) (resp *types.DeviceMsgSendLogIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.SendLogIndex(l.ctx, utils.Copy[dm.SendLogIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.HubLogIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgSendLogInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		var user *types.UserCore
		if req.WithUser {
			if v.UserID <= def.RootNode {
				user = &types.UserCore{
					UserID:   v.UserID,
					UserName: v.Account,
				}
			} else {
				ui, err := l.svcCtx.UserC.GetData(l.ctx, v.UserID)
				if err == nil {
					user = utils.Copy[types.UserCore](ui)
				}
			}
		}
		info = append(info, &types.DeviceMsgSendLogInfo{
			Timestamp:  v.Timestamp,
			Account:    v.Account,
			UserID:     v.UserID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
			Action:     v.Action,
			DataID:     v.DataID,
			TraceID:    v.TraceID,
			Content:    v.Content,
			ResultCode: v.ResultCode,
			User:       user,
		})
	}
	return &types.DeviceMsgSendLogIndexResp{List: info, PageResp: logic.ToPageResp(req.Page, dmResp.Total)}, err
}
