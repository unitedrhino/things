package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) UserDelete(in *sys.UserDeleteReq) (*sys.Response, error) {
	err := l.svcCtx.UserInfoModel.Delete(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("%s.Delete uid=%d err=%+v", utils.FuncName(), in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}

	l.Infof("%s.delete uid=%v", utils.FuncName(), in.Uid)

	return &sys.Response{}, nil
}
