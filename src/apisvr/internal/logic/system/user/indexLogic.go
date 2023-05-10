package user

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.UserIndexReq) (resp *types.UserIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	var page sys.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.UserRpc.UserIndex(l.ctx, &sys.UserIndexReq{
		Page:     &page,
		UserName: req.UserName,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}

	var userInfo []*types.UserInfo
	var total int64
	total = info.Total

	userInfo = make([]*types.UserInfo, 0, len(userInfo))
	for _, i := range info.Info {
		userInfo = append(userInfo, UserInfoToApi(i))
	}

	return &types.UserIndexResp{userInfo, total}, nil
}
