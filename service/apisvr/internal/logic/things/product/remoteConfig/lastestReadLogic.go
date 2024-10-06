package remoteConfig

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LastestReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLastestReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LastestReadLogic {
	return &LastestReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LastestReadLogic) LastestRead(req *types.ProductRemoteConfigLastestReadReq) (resp *types.ProductRemoteConfigLastestReadResp, err error) {
	res, err := l.svcCtx.RemoteConfig.RemoteConfigLastRead(l.ctx, &dm.RemoteConfigLastReadReq{
		ProductID: req.ProductID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.RemoteConfigLastRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	return &types.ProductRemoteConfigLastestReadResp{types.ProductRemoteConfig{
		ID:         res.Info.Id,
		Content:    res.Info.Content,
		CreateTime: cast.ToString(res.Info.CreatedTime),
	}}, nil
}
