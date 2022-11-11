package remoteconfiglogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigLastReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoteConfigLastReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigLastReadLogic {
	return &RemoteConfigLastReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoteConfigLastReadLogic) RemoteConfigLastRead(in *dm.RemoteConfigLastReadReq) (*dm.RemoteConfigLastReadResp, error) {
	res, err := l.svcCtx.RemoteConfigDB.GetLastRecord(l.ctx, &mysql.RemoteConfigFilter{
		ProductID: in.ProductID,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
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
		CreatedTime: res.CreatedTime,
	}}, nil
}
