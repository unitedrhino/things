package remoteConfig

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.ProductRemoteConfigIndexReq) (resp *types.ProductRemoteConfigIndexResp, err error) {
	res, err := l.svcCtx.RemoteConfig.RemoteConfigIndex(l.ctx, &dm.RemoteConfigIndexReq{
		Page:      &dm.PageInfo{Page: req.Page.Page, Size: req.Page.Size},
		ProductID: req.ProductID,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.RemoteConfigIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	list := make([]*types.ProductRemoteConfig, 0, len(res.List))
	for _, v := range res.List {
		list = append(list, &types.ProductRemoteConfig{
			ID:         v.Id,
			Content:    v.Content,
			CreateTime: cast.ToString(v.CreatedTime),
		})
	}
	return &types.ProductRemoteConfigIndexResp{
		List:  list,
		Total: res.Total,
	}, nil
}
