package user

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/usersvr/user"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Register2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) Register2Logic {
	return Register2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//注册完成后就需要填写用户信息,填写完成后才算注册成功(目前只有手机号注册登录需要走这步)
func (l *Register2Logic) Register2(req types.Register2Req) error {
	l.Infof("Register2|req=%+v", req)
	token, err := users.ParseToken(req.Token, l.svcCtx.Config.Rej.AccessSecret)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("parseToken failure|token=%s|err=%+v", token, er)
		return er
	}
	if token.Uid != req.Uid {
		l.Errorf("uid is invalid")
		return errors.UidNotRight
	}
	resp, er := l.svcCtx.UserRpc.Register2(l.ctx, &user.Register2Req{
		Password: req.Password,
		Info: &user.UserInfo{
			Uid:        req.Uid,
			UserName:   req.UserName,
			NickName:   req.NickName,
			InviterUid: req.InviterUid,
			InviterId:  req.InviterId,
			Sex:        req.Sex,
			City:       req.City,
			Country:    req.Country,
			Province:   req.Province,
			Language:   req.Language,
			HeadImgUrl: req.HeadImgUrl,
		},
	})
	if er != nil {
		err := errors.Fmt(er)
		l.Errorf("[%s]|rpc.RegisterCore|req=%v|err=%+v", utils.FuncName(), req, err)
		return err
	}
	if resp == nil {
		l.Errorf("%s|rpc.RegisterCore|return nil|req=%+v", utils.FuncName(), req)
		return errors.System.AddDetail("register core rpc return nil")
	}
	return nil
}
