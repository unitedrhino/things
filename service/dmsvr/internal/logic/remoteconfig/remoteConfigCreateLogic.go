package remoteconfiglogic

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PrcDB *relationDB.ProductRemoteConfigRepo
}

func NewRemoteConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigCreateLogic {
	return &RemoteConfigCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PrcDB:  relationDB.NewProductRemoteConfigRepo(ctx),
	}
}

func (l *RemoteConfigCreateLogic) RemoteConfigCreate(in *dm.RemoteConfigCreateReq) (*dm.Response, error) {
	err := l.PrcDB.Insert(l.ctx, &relationDB.DmProductRemoteConfig{
		ProductID: in.ProductID,
		Content:   in.Content,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &dm.Response{}, nil
}
