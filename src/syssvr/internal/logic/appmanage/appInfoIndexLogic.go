package appmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoIndexLogic {
	return &AppInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoIndexLogic) AppInfoIndex(in *sys.AppInfoIndexReq) (*sys.AppInfoIndexResp, error) {
	ret, err := relationDB.NewAppInfoRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.AppInfoFilter{Code: in.Code, Name: in.Name}, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &sys.AppInfoIndexResp{List: ToAppInfosPb(ret), Total: int64(len(ret))}, nil
}
