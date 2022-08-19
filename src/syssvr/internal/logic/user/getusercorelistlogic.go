package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCoreListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCoreListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCoreListLogic {
	return &GetUserCoreListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCoreListLogic) GetUserCoreList(in *sys.GetUserCoreListReq) (*sys.GetUserCoreListResp, error) {
	l.Infof("GetUserCoreList|req=%+v", in)
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	ucs, total, err := l.svcCtx.UserModel.GetUserCoreList(page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	info := make([]*sys.UserCore, 0, len(ucs))
	for _, uc := range ucs {
		info = append(info, UserCoreToPb(uc))
	}
	return &sys.GetUserCoreListResp{
		Info:  info,
		Total: total,
	}, nil
}
