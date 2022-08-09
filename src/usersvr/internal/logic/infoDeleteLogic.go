package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/pb/user"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type InfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoDeleteLogic {
	return &InfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InfoDeleteLogic) InfoDelete(in *user.UserInfoDeleteReq) (*user.Response, error) {
	l.Infof("ModifyUserInfo|req=%+v", in)
	var err error
	//todo : 后续补充先通过uid查找两张表，是否同时存在记录，并且用户名使用一致，如果条件都满足，再对两表记录进行删除，目前的实现比较简单
	err = l.svcCtx.UserInfoModel.Delete(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("UserInfoModel|Delete|uid=%d|err=%+v", in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}

	err = l.svcCtx.UserCoreModel.Delete(l.ctx, cast.ToInt64(in.Uid))
	if err != nil {
		l.Errorf("UserCoreModel|Delete|uid=%d|err=%+v", in.Uid, err)
		return nil, errors.Database.AddDetail(err)
	}

	l.Infof("InfoDelete|delete usersvr uid= %s", in.Uid)

	return &user.Response{}, nil
}
