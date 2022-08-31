package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
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

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *sys.UserDeleteReq) (*sys.Response, error) {
	err := l.svcCtx.UserInfoModel.Delete(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("UserInfoModel|Delete|uid=%d|err=%+v", in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}

	l.Infof("InfoDelete|delete usersvr uid= %s", in.Uid)

	return &sys.Response{}, nil
}
