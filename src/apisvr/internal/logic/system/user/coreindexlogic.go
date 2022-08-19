package user

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type CoreIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoreIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoreIndexLogic {
	return &CoreIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoreIndexLogic) CoreIndex(req *types.UserCoreIndexReq) (resp *types.UserCoreIndexResp, err error) {
	l.Infof("UserCoreList|req=%+v", req)
	var page sys.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.UserRpc.GetUserCoreList(l.ctx, &sys.GetUserCoreListReq{
		Page: &page,
	})
	if err != nil {
		return nil, err
	}

	var usercore []*types.UserCore
	var total int64
	total = info.Total

	usercore = make([]*types.UserCore, 0, len(usercore))
	for _, i := range info.Info {
		usercore = append(usercore, UserCoreToApi(i))
	}

	return &types.UserCoreIndexResp{usercore, total}, nil
}
