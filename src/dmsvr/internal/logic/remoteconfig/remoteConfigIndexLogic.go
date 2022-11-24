package remoteconfiglogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoteConfigIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigIndexLogic {
	return &RemoteConfigIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoteConfigIndexLogic) RemoteConfigIndex(in *dm.RemoteConfigIndexReq) (*dm.RemoteConfigIndexResp, error) {
	resp, total, err := l.svcCtx.RemoteConfigDB.Index(l.ctx, &mysql.RemoteConfigFilter{
		Page:      &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		ProductID: in.ProductID,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	list := make([]*dm.ProductRemoteConfig, 0, len(resp))
	for _, v := range resp {
		list = append(list, &dm.ProductRemoteConfig{
			Id:          v.ID,
			ProductID:   v.ProductID,
			Content:     v.Content,
			CreatedTime: v.CreatedTime,
		})
	}

	return &dm.RemoteConfigIndexResp{List: list, Total: total}, nil
}
