package remoteconfiglogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoteConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigCreateLogic {
	return &RemoteConfigCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoteConfigCreateLogic) RemoteConfigCreate(in *dm.RemoteConfigCreateReq) (*dm.Response, error) {
	_, err := l.svcCtx.RemoteConfigInfo.Insert(l.ctx, &mysql.DmProductRemoteConfig{
		ProductID: in.ProductID,
		Content:   in.Content,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &dm.Response{}, nil
}
