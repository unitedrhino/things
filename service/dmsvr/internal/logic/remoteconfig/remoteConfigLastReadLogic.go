package remoteconfiglogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigLastReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PrcDB *relationDB.ProductRemoteConfigRepo
}

func NewRemoteConfigLastReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigLastReadLogic {
	return &RemoteConfigLastReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PrcDB:  relationDB.NewProductRemoteConfigRepo(ctx),
	}
}

func (l *RemoteConfigLastReadLogic) RemoteConfigLastRead(in *dm.RemoteConfigLastReadReq) (*dm.RemoteConfigLastReadResp, error) {
	res, err := l.PrcDB.FindOneByFilter(l.ctx, relationDB.RemoteConfigFilter{
		ProductID: in.ProductID,
	})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}
	if res == nil {
		return &dm.RemoteConfigLastReadResp{Info: &dm.ProductRemoteConfig{
			Id:          0,
			ProductID:   in.ProductID,
			Content:     "",
			CreatedTime: 0,
		}}, nil
	}

	return &dm.RemoteConfigLastReadResp{Info: &dm.ProductRemoteConfig{
		Id:          res.ID,
		ProductID:   res.ProductID,
		Content:     res.Content,
		CreatedTime: res.CreatedTime.Unix(),
	}}, nil
}
