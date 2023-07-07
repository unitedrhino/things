package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *IndexLogic) UserIndex(in *sys.UserIndexReq) (*sys.UserIndexResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	f := relationDB.UserInfoFilter{
		UserName: in.UserName,
		Phone:    in.Phone,
		Email:    in.Email,
	}
	ucs, err := l.UiDB.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := l.UiDB.CountByFilter(l.ctx, f)
	info := make([]*sys.UserInfo, 0, len(ucs))
	for _, uc := range ucs {
		info = append(info, UserInfoToPb(uc))
	}
	if err != nil {
		return nil, err
	}
	return &sys.UserIndexResp{
		Info:  info,
		Total: total,
	}, nil

}
