package logic

import (
	"context"
	"github.com/go-things/things/shared/def"
	"github.com/go-things/things/shared/errors"
	"github.com/jinzhu/copier"

	"github.com/go-things/things/src/usersvr/internal/svc"
	"github.com/go-things/things/src/usersvr/user"

	"github.com/tal-tech/go-zero/core/logx"
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

func (l *GetUserCoreListLogic) GetUserCoreList(in *user.GetUserCoreListReq) (*user.GetUserCoreListResp, error) {
	l.Infof("GetUserCoreList|req=%+v", in)
	page := def.PageInfo{}
	copier.Copy(&page, in.Page)
	ucs, total, err := l.svcCtx.UserModel.GetUserCoreList(page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	info := make([]*user.UserCore, 0, len(ucs))
	for _, uc := range ucs {
		info = append(info, UserCoreToPb(uc))
	}
	return &user.GetUserCoreListResp{
		Info:  info,
		Total: total,
	}, nil
}
