package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/usersvr/internal/svc"
	"gitee.com/godLei6/things/src/usersvr/model"
	"gitee.com/godLei6/things/src/usersvr/user"
	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserCoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCoreLogic {
	return &GetUserCoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCoreLogic) GetUserCore(in *user.GetUserCoreReq) (*user.GetUserCoreResp, error) {
	l.Infof("GetUserCore|req=%+v", in)
	uc, err := l.svcCtx.UserCoreModel.FindOne(in.Uid)
	switch err {
	case nil:
		return &user.GetUserCoreResp{Info: UserCoreToPb(uc)}, nil
	case model.ErrNotFound:
		return nil, errors.UidNotRight
	default:
		l.Errorf("GetUserCore|req=%#v|err=%+v", in, err)
		return nil, errors.Database.AddDetail(err)
	}
}
