package info

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/user"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"

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

func (l *IndexLogic) Index(req *types.UserInfoIndexReq) (resp *types.UserInfoIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	info, err := l.svcCtx.UserRpc.UserInfoIndex(l.ctx, &sys.UserInfoIndexReq{
		Page:     logic.ToSysPageRpc(req.Page),
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
	for _, i := range info.List {
		userInfo = append(userInfo, user.UserInfoToApi(i, nil, nil))
	}

	return &types.UserInfoIndexResp{userInfo, total}, nil
}
