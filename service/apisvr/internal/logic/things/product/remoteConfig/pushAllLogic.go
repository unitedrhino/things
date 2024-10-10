package remoteConfig

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPushAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushAllLogic {
	return &PushAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushAllLogic) PushAll(req *types.ProductRemoteConfigPushAllReq) error {
	_, err := l.svcCtx.RemoteConfig.RemoteConfigPushAll(l.ctx, &dm.RemoteConfigPushAllReq{
		ProductID: req.ProductID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.RemoteConfigPushAll req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}

	return nil
}
