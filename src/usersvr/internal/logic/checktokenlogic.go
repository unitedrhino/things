package logic

import (
	"context"
	"github.com/i-Things/things/shared/users"
	"github.com/i-Things/things/shared/utils"
	"time"

	"github.com/i-Things/things/src/usersvr/internal/svc"
	"github.com/i-Things/things/src/usersvr/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckTokenLogic) CheckToken(in *user.CheckTokenReq) (*user.CheckTokenResp, error) {
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(l.ctx, p)
		}
	}()
	l.Infof("CheckToken|req=%+v", in)
	jwt, err := users.ParseToken(in.Token, l.svcCtx.Config.UserToken.AccessSecret)
	if err != nil {
		l.Errorf("CheckToken|parse token fail|err=%s", err.Error())
		return nil, err
	}
	var token string
	if (jwt.ExpiresAt-time.Now().Unix())*2 < l.svcCtx.Config.UserToken.AccessExpire {
		token, _ = users.RefreshToken(in.Token, l.svcCtx.Config.UserToken.AccessSecret)
	}
	l.Infof("CheckToken|uid=%d", jwt.Uid)
	return &user.CheckTokenResp{
		Token: token,
		Uid:   jwt.Uid,
	}, nil
}
